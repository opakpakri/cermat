package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cermat-api/internal/middleware"
	"cermat-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DealHandler struct {
	DB *pgxpool.Pool
}

type dealRequest struct {
	Title             string  `json:"title"`
	ContactID         int     `json:"contact_id"`
	CatalogItemID     *int    `json:"catalog_item_id"`
	Stage             string  `json:"stage"`
	DealValue         float64 `json:"deal_value"`
	PaymentPlan       *string `json:"payment_plan"`
	ExpectedCloseDate *string `json:"expected_close_date"`
	Notes             *string `json:"notes"`
}

func (h *DealHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	stage := strings.ToUpper(strings.TrimSpace(c.Query("stage")))
	search := strings.TrimSpace(c.Query("search"))
	leadSource := strings.ToUpper(strings.TrimSpace(c.Query("lead_source")))

	query := `
		SELECT d.id, d.user_id, d.deal_code, d.title, d.contact_id, c.name, c.phone, c.lead_source,
		       d.catalog_item_id, ci.name, ci.code, d.stage, d.deal_value, d.payment_plan,
		       d.expected_close_date, d.notes, d.created_at, d.updated_at
		FROM deals d
		JOIN contacts c ON c.id = d.contact_id AND c.user_id = d.user_id
		LEFT JOIN catalog_items ci ON ci.id = d.catalog_item_id AND ci.user_id = d.user_id
		WHERE d.user_id = $1
	`
	args := []interface{}{userID}
	idx := 2

	if stage != "" {
		query += fmt.Sprintf(" AND d.stage = $%d", idx)
		args = append(args, stage)
		idx++
	}
	if leadSource != "" {
		query += fmt.Sprintf(" AND c.lead_source = $%d", idx)
		args = append(args, leadSource)
		idx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (d.title ILIKE $%d OR c.name ILIKE $%d OR d.deal_code ILIKE $%d)", idx, idx, idx)
		args = append(args, "%"+search+"%")
		idx++
	}
	query += " ORDER BY d.updated_at DESC"

	rows, err := h.DB.Query(c.Context(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat deals"})
	}
	defer rows.Close()

	deals := make([]models.Deal, 0)
	for rows.Next() {
		var dl models.Deal
		if err := rows.Scan(
			&dl.ID, &dl.UserID, &dl.DealCode, &dl.Title, &dl.ContactID, &dl.ContactName, &dl.ContactPhone, &dl.LeadSource,
			&dl.CatalogItemID, &dl.CatalogItemName, &dl.CatalogItemCode, &dl.Stage, &dl.DealValue, &dl.PaymentPlan,
			&dl.ExpectedCloseDate, &dl.Notes, &dl.CreatedAt, &dl.UpdatedAt,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data deal"})
		}
		deals = append(deals, dl)
	}

	return c.JSON(deals)
}

func (h *DealHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req dealRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Stage = strings.ToUpper(strings.TrimSpace(req.Stage))
	if req.Title == "" || req.ContactID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Judul deal dan kontak wajib diisi"})
	}
	if req.Stage == "" {
		req.Stage = "NEW"
	}

	// Cross-Tenant Reference Validation 1: contact_id must belong to userID
	var contactExists bool
	err := h.DB.QueryRow(c.Context(), `SELECT EXISTS(SELECT 1 FROM contacts WHERE id = $1 AND user_id = $2)`, req.ContactID, userID).Scan(&contactExists)
	if err != nil || !contactExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kontak tidak ditemukan"})
	}

	// Cross-Tenant Reference Validation 2: catalog_item_id (if provided) must belong to userID
	if req.CatalogItemID != nil && *req.CatalogItemID > 0 {
		var itemPrice float64
		err = h.DB.QueryRow(c.Context(), `SELECT price FROM catalog_items WHERE id = $1 AND user_id = $2`, *req.CatalogItemID, userID).Scan(&itemPrice)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item katalog tidak ditemukan"})
		}
		if req.DealValue == 0 {
			req.DealValue = itemPrice
		}
	} else {
		req.CatalogItemID = nil
	}

	var closeDate *time.Time
	if req.ExpectedCloseDate != nil && strings.TrimSpace(*req.ExpectedCloseDate) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*req.ExpectedCloseDate))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal perkiraan closing tidak valid (YYYY-MM-DD)"})
		}
		closeDate = &parsed
	}

	// Generate unique deal_code per tenant
	var count int
	_ = h.DB.QueryRow(c.Context(), `SELECT COUNT(*) FROM deals WHERE user_id = $1`, userID).Scan(&count)
	dealCode := fmt.Sprintf("CRM-%s-%03d", time.Now().Format("20060102"), count+1)

	var dealID int
	err = h.DB.QueryRow(c.Context(), `
		INSERT INTO deals (user_id, deal_code, title, contact_id, catalog_item_id, stage, deal_value, payment_plan, expected_close_date, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, userID, dealCode, req.Title, req.ContactID, req.CatalogItemID, req.Stage, req.DealValue, req.PaymentPlan, closeDate, req.Notes).Scan(&dealID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat deal"})
	}

	// Auto-update catalog item status if applicable
	if req.CatalogItemID != nil {
		h.syncCatalogStatus(c, userID, *req.CatalogItemID, req.Stage)
	}

	return h.fetchAndReturnDeal(c, dealID, userID, fiber.StatusCreated)
}

func (h *DealHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Deal tidak ditemukan"})
	}

	var dl models.Deal
	err = h.DB.QueryRow(c.Context(), `
		SELECT d.id, d.user_id, d.deal_code, d.title, d.contact_id, c.name, c.phone, c.lead_source,
		       d.catalog_item_id, ci.name, ci.code, d.stage, d.deal_value, d.payment_plan,
		       d.expected_close_date, d.notes, d.created_at, d.updated_at
		FROM deals d
		JOIN contacts c ON c.id = d.contact_id AND c.user_id = d.user_id
		LEFT JOIN catalog_items ci ON ci.id = d.catalog_item_id AND ci.user_id = d.user_id
		WHERE d.id = $1 AND d.user_id = $2
	`, id, userID).Scan(
		&dl.ID, &dl.UserID, &dl.DealCode, &dl.Title, &dl.ContactID, &dl.ContactName, &dl.ContactPhone, &dl.LeadSource,
		&dl.CatalogItemID, &dl.CatalogItemName, &dl.CatalogItemCode, &dl.Stage, &dl.DealValue, &dl.PaymentPlan,
		&dl.ExpectedCloseDate, &dl.Notes, &dl.CreatedAt, &dl.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Deal tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat detail deal"})
	}

	activities := make([]models.DealActivity, 0)
	actRows, err := h.DB.Query(c.Context(), `
		SELECT a.id, a.user_id, a.deal_id, d.title, a.contact_id, c.name, c.phone,
		       a.activity_type, a.title, a.description, a.scheduled_at, a.is_completed, a.created_at
		FROM deal_activities a
		JOIN contacts c ON c.id = a.contact_id AND c.user_id = a.user_id
		LEFT JOIN deals d ON d.id = a.deal_id AND d.user_id = a.user_id
		WHERE a.deal_id = $1 AND a.user_id = $2
		ORDER BY a.created_at DESC
	`, id, userID)
	if err == nil {
		defer actRows.Close()
		for actRows.Next() {
			var act models.DealActivity
			if scanErr := actRows.Scan(
				&act.ID, &act.UserID, &act.DealID, &act.DealTitle, &act.ContactID, &act.ContactName, &act.ContactPhone,
				&act.ActivityType, &act.Title, &act.Description, &act.ScheduledAt, &act.IsCompleted, &act.CreatedAt,
			); scanErr == nil {
				activities = append(activities, act)
			}
		}
	}

	return c.JSON(fiber.Map{
		"deal":       dl,
		"activities": activities,
	})
}

func (h *DealHandler) UpdateStage(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Deal tidak ditemukan"})
	}

	var req struct {
		Stage string `json:"stage"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}
	req.Stage = strings.ToUpper(strings.TrimSpace(req.Stage))
	validStages := map[string]bool{
		"NEW": true, "QUALIFIED": true, "PROPOSAL_SURVEY": true,
		"NEGOTIATION_BOOKING": true, "WON": true, "LOST": true,
	}
	if !validStages[req.Stage] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tahapan pipeline tidak valid"})
	}

	var catalogItemID *int
	err = h.DB.QueryRow(c.Context(), `
		UPDATE deals
		SET stage = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND user_id = $3
		RETURNING catalog_item_id
	`, req.Stage, id, userID).Scan(&catalogItemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Deal tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengubah tahapan deal"})
	}

	if catalogItemID != nil {
		h.syncCatalogStatus(c, userID, *catalogItemID, req.Stage)
	}

	return h.fetchAndReturnDeal(c, id, userID, fiber.StatusOK)
}

func (h *DealHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Deal tidak ditemukan"})
	}

	var req dealRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Stage = strings.ToUpper(strings.TrimSpace(req.Stage))
	if req.Title == "" || req.ContactID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Judul deal dan kontak wajib diisi"})
	}

	// Validate contact ownership
	var contactExists bool
	err = h.DB.QueryRow(c.Context(), `SELECT EXISTS(SELECT 1 FROM contacts WHERE id = $1 AND user_id = $2)`, req.ContactID, userID).Scan(&contactExists)
	if err != nil || !contactExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kontak tidak ditemukan"})
	}

	if req.CatalogItemID != nil && *req.CatalogItemID > 0 {
		var itemExists bool
		err = h.DB.QueryRow(c.Context(), `SELECT EXISTS(SELECT 1 FROM catalog_items WHERE id = $1 AND user_id = $2)`, *req.CatalogItemID, userID).Scan(&itemExists)
		if err != nil || !itemExists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item katalog tidak ditemukan"})
		}
	} else {
		req.CatalogItemID = nil
	}

	var closeDate *time.Time
	if req.ExpectedCloseDate != nil && strings.TrimSpace(*req.ExpectedCloseDate) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*req.ExpectedCloseDate))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal tidak valid"})
		}
		closeDate = &parsed
	}

	cmd, err := h.DB.Exec(c.Context(), `
		UPDATE deals
		SET title = $1, contact_id = $2, catalog_item_id = $3, stage = $4,
		    deal_value = $5, payment_plan = $6, expected_close_date = $7, notes = $8,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $9 AND user_id = $10
	`, req.Title, req.ContactID, req.CatalogItemID, req.Stage, req.DealValue, req.PaymentPlan, closeDate, req.Notes, id, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui deal"})
	}
	if cmd.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Deal tidak ditemukan"})
	}

	if req.CatalogItemID != nil {
		h.syncCatalogStatus(c, userID, *req.CatalogItemID, req.Stage)
	}

	return h.fetchAndReturnDeal(c, id, userID, fiber.StatusOK)
}

func (h *DealHandler) syncCatalogStatus(c *fiber.Ctx, userID int, catalogItemID int, stage string) {
	var newStatus string
	switch stage {
	case "NEGOTIATION_BOOKING":
		newStatus = "BOOKED"
	case "WON":
		newStatus = "SOLD"
	default:
		return
	}
	_, _ = h.DB.Exec(c.Context(), `
		UPDATE catalog_items SET status = $1 WHERE id = $2 AND user_id = $3
	`, newStatus, catalogItemID, userID)
}

func (h *DealHandler) fetchAndReturnDeal(c *fiber.Ctx, dealID int, userID int, statusCode int) error {
	var dl models.Deal
	err := h.DB.QueryRow(c.Context(), `
		SELECT d.id, d.user_id, d.deal_code, d.title, d.contact_id, c.name, c.phone, c.lead_source,
		       d.catalog_item_id, ci.name, ci.code, d.stage, d.deal_value, d.payment_plan,
		       d.expected_close_date, d.notes, d.created_at, d.updated_at
		FROM deals d
		JOIN contacts c ON c.id = d.contact_id AND c.user_id = d.user_id
		LEFT JOIN catalog_items ci ON ci.id = d.catalog_item_id AND ci.user_id = d.user_id
		WHERE d.id = $1 AND d.user_id = $2
	`, dealID, userID).Scan(
		&dl.ID, &dl.UserID, &dl.DealCode, &dl.Title, &dl.ContactID, &dl.ContactName, &dl.ContactPhone, &dl.LeadSource,
		&dl.CatalogItemID, &dl.CatalogItemName, &dl.CatalogItemCode, &dl.Stage, &dl.DealValue, &dl.PaymentPlan,
		&dl.ExpectedCloseDate, &dl.Notes, &dl.CreatedAt, &dl.UpdatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat deal"})
	}
	return c.Status(statusCode).JSON(dl)
}
