package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"cermat-api/internal/middleware"
	"cermat-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactHandler struct {
	DB *pgxpool.Pool
}

type contactRequest struct {
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	Email       *string `json:"email"`
	CompanyName *string `json:"company_name"`
	Address     *string `json:"address"`
	LeadSource  string  `json:"lead_source"`
	Status      string  `json:"status"`
	Notes       *string `json:"notes"`
}

func (h *ContactHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	status := strings.ToUpper(strings.TrimSpace(c.Query("status")))
	search := strings.TrimSpace(c.Query("search"))

	query := `
		SELECT id, user_id, name, phone, email, company_name, address, lead_source, status, notes, created_at
		FROM contacts
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	idx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, status)
		idx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR phone ILIKE $%d OR COALESCE(company_name, '') ILIKE $%d)", idx, idx, idx)
		args = append(args, "%"+search+"%")
		idx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := h.DB.Query(c.Context(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat kontak"})
	}
	defer rows.Close()

	contacts := make([]models.Contact, 0)
	for rows.Next() {
		var ct models.Contact
		if err := rows.Scan(
			&ct.ID, &ct.UserID, &ct.Name, &ct.Phone, &ct.Email, &ct.CompanyName,
			&ct.Address, &ct.LeadSource, &ct.Status, &ct.Notes, &ct.CreatedAt,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data kontak"})
		}
		contacts = append(contacts, ct)
	}

	return c.JSON(contacts)
}

func (h *ContactHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req contactRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.LeadSource = strings.ToUpper(strings.TrimSpace(req.LeadSource))
	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))

	if req.Name == "" || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama dan no. telepon wajib diisi"})
	}
	if req.LeadSource == "" {
		req.LeadSource = "WHATSAPP"
	}
	if req.Status == "" {
		req.Status = "LEAD"
	}

	var ct models.Contact
	err := h.DB.QueryRow(c.Context(), `
		INSERT INTO contacts (user_id, name, phone, email, company_name, address, lead_source, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, name, phone, email, company_name, address, lead_source, status, notes, created_at
	`, userID, req.Name, req.Phone, req.Email, req.CompanyName, req.Address, req.LeadSource, req.Status, req.Notes).Scan(
		&ct.ID, &ct.UserID, &ct.Name, &ct.Phone, &ct.Email, &ct.CompanyName,
		&ct.Address, &ct.LeadSource, &ct.Status, &ct.Notes, &ct.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Nomor telepon sudah terdaftar di kontak Anda"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menambah kontak"})
	}

	return c.Status(fiber.StatusCreated).JSON(ct)
}

func (h *ContactHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kontak tidak ditemukan"})
	}

	var ct models.Contact
	err = h.DB.QueryRow(c.Context(), `
		SELECT id, user_id, name, phone, email, company_name, address, lead_source, status, notes, created_at
		FROM contacts
		WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(
		&ct.ID, &ct.UserID, &ct.Name, &ct.Phone, &ct.Email, &ct.CompanyName,
		&ct.Address, &ct.LeadSource, &ct.Status, &ct.Notes, &ct.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kontak tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat detail kontak"})
	}

	deals := make([]models.Deal, 0)
	rows, err := h.DB.Query(c.Context(), `
		SELECT d.id, d.user_id, d.deal_code, d.title, d.contact_id, c.name, c.phone, c.lead_source,
		       d.catalog_item_id, ci.name, ci.code, d.stage, d.deal_value, d.payment_plan,
		       d.expected_close_date, d.notes, d.created_at, d.updated_at
		FROM deals d
		JOIN contacts c ON c.id = d.contact_id AND c.user_id = d.user_id
		LEFT JOIN catalog_items ci ON ci.id = d.catalog_item_id AND ci.user_id = d.user_id
		WHERE d.contact_id = $1 AND d.user_id = $2
		ORDER BY d.created_at DESC
	`, id, userID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var dl models.Deal
			if scanErr := rows.Scan(
				&dl.ID, &dl.UserID, &dl.DealCode, &dl.Title, &dl.ContactID, &dl.ContactName, &dl.ContactPhone, &dl.LeadSource,
				&dl.CatalogItemID, &dl.CatalogItemName, &dl.CatalogItemCode, &dl.Stage, &dl.DealValue, &dl.PaymentPlan,
				&dl.ExpectedCloseDate, &dl.Notes, &dl.CreatedAt, &dl.UpdatedAt,
			); scanErr == nil {
				deals = append(deals, dl)
			}
		}
	}

	return c.JSON(fiber.Map{
		"contact": ct,
		"deals":   deals,
	})
}

func (h *ContactHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kontak tidak ditemukan"})
	}

	var req contactRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.LeadSource = strings.ToUpper(strings.TrimSpace(req.LeadSource))
	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))

	if req.Name == "" || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama dan no. telepon wajib diisi"})
	}

	var ct models.Contact
	err = h.DB.QueryRow(c.Context(), `
		UPDATE contacts
		SET name = $1, phone = $2, email = $3, company_name = $4, address = $5,
		    lead_source = $6, status = $7, notes = $8
		WHERE id = $9 AND user_id = $10
		RETURNING id, user_id, name, phone, email, company_name, address, lead_source, status, notes, created_at
	`, req.Name, req.Phone, req.Email, req.CompanyName, req.Address, req.LeadSource, req.Status, req.Notes, id, userID).Scan(
		&ct.ID, &ct.UserID, &ct.Name, &ct.Phone, &ct.Email, &ct.CompanyName,
		&ct.Address, &ct.LeadSource, &ct.Status, &ct.Notes, &ct.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kontak tidak ditemukan"})
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Nomor telepon sudah digunakan kontak lain"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui kontak"})
	}

	return c.JSON(ct)
}

func (h *ContactHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kontak tidak ditemukan"})
	}

	cmd, err := h.DB.Exec(c.Context(), `DELETE FROM contacts WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Kontak tidak dapat dihapus karena masih terhubung dengan deal atau tiket",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus kontak"})
	}
	if cmd.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kontak tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Kontak berhasil dihapus"})
}
