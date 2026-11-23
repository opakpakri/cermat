package handlers

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"cermat-api/internal/middleware"
	"cermat-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityHandler struct {
	DB *pgxpool.Pool
}

type activityRequest struct {
	DealID       *int    `json:"deal_id"`
	ContactID    int     `json:"contact_id"`
	ActivityType string  `json:"activity_type"`
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	ScheduledAt  *string `json:"scheduled_at"`
}

func (h *ActivityHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	completedParam := strings.ToLower(strings.TrimSpace(c.Query("completed")))

	query := `
		SELECT a.id, a.user_id, a.deal_id, d.title, a.contact_id, c.name, c.phone,
		       a.activity_type, a.title, a.description, a.scheduled_at, a.is_completed, a.created_at
		FROM deal_activities a
		JOIN contacts c ON c.id = a.contact_id AND c.user_id = a.user_id
		LEFT JOIN deals d ON d.id = a.deal_id AND d.user_id = a.user_id
		WHERE a.user_id = $1
	`
	args := []interface{}{userID}
	if completedParam == "false" {
		query += " AND a.is_completed = FALSE"
	} else if completedParam == "true" {
		query += " AND a.is_completed = TRUE"
	}
	query += " ORDER BY a.is_completed ASC, COALESCE(a.scheduled_at, a.created_at) ASC"

	rows, err := h.DB.Query(c.Context(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat aktivitas"})
	}
	defer rows.Close()

	activities := make([]models.DealActivity, 0)
	for rows.Next() {
		var act models.DealActivity
		if err := rows.Scan(
			&act.ID, &act.UserID, &act.DealID, &act.DealTitle, &act.ContactID, &act.ContactName, &act.ContactPhone,
			&act.ActivityType, &act.Title, &act.Description, &act.ScheduledAt, &act.IsCompleted, &act.CreatedAt,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data aktivitas"})
		}
		activities = append(activities, act)
	}

	return c.JSON(activities)
}

func (h *ActivityHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req activityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Title = strings.TrimSpace(req.Title)
	req.ActivityType = strings.ToUpper(strings.TrimSpace(req.ActivityType))
	if req.Title == "" || req.ContactID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Judul aktivitas dan kontak wajib diisi"})
	}
	if req.ActivityType == "" {
		req.ActivityType = "CALL"
	}

	// Cross-tenant validation: contact_id
	var contactExists bool
	err := h.DB.QueryRow(c.Context(), `SELECT EXISTS(SELECT 1 FROM contacts WHERE id = $1 AND user_id = $2)`, req.ContactID, userID).Scan(&contactExists)
	if err != nil || !contactExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kontak tidak ditemukan"})
	}

	// Cross-tenant validation: deal_id (if provided)
	if req.DealID != nil && *req.DealID > 0 {
		var dealExists bool
		err = h.DB.QueryRow(c.Context(), `SELECT EXISTS(SELECT 1 FROM deals WHERE id = $1 AND user_id = $2)`, *req.DealID, userID).Scan(&dealExists)
		if err != nil || !dealExists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Deal tidak ditemukan"})
		}
	} else {
		req.DealID = nil
	}

	var scheduledAt *time.Time
	if req.ScheduledAt != nil && strings.TrimSpace(*req.ScheduledAt) != "" {
		raw := strings.TrimSpace(*req.ScheduledAt)
		parsed, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			parsed, err = time.Parse("2006-01-02T15:04", raw)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format waktu jadwal tidak valid"})
			}
		}
		scheduledAt = &parsed
	}

	var actID int
	err = h.DB.QueryRow(c.Context(), `
		INSERT INTO deal_activities (user_id, deal_id, contact_id, activity_type, title, description, scheduled_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, userID, req.DealID, req.ContactID, req.ActivityType, req.Title, req.Description, scheduledAt).Scan(&actID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menambah aktivitas"})
	}

	var act models.DealActivity
	err = h.DB.QueryRow(c.Context(), `
		SELECT a.id, a.user_id, a.deal_id, d.title, a.contact_id, c.name, c.phone,
		       a.activity_type, a.title, a.description, a.scheduled_at, a.is_completed, a.created_at
		FROM deal_activities a
		JOIN contacts c ON c.id = a.contact_id AND c.user_id = a.user_id
		LEFT JOIN deals d ON d.id = a.deal_id AND d.user_id = a.user_id
		WHERE a.id = $1 AND a.user_id = $2
	`, actID, userID).Scan(
		&act.ID, &act.UserID, &act.DealID, &act.DealTitle, &act.ContactID, &act.ContactName, &act.ContactPhone,
		&act.ActivityType, &act.Title, &act.Description, &act.ScheduledAt, &act.IsCompleted, &act.CreatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat aktivitas"})
	}

	return c.Status(fiber.StatusCreated).JSON(act)
}

func (h *ActivityHandler) Toggle(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Aktivitas tidak ditemukan"})
	}

	var currentCompleted bool
	err = h.DB.QueryRow(c.Context(), `
		SELECT is_completed FROM deal_activities WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&currentCompleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Aktivitas tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memeriksa status aktivitas"})
	}

	if currentCompleted {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Aktivitas yang sudah selesai tidak dapat diubah kembali",
		})
	}

	var isCompleted bool
	err = h.DB.QueryRow(c.Context(), `
		UPDATE deal_activities
		SET is_completed = TRUE
		WHERE id = $1 AND user_id = $2
		RETURNING is_completed
	`, id, userID).Scan(&isCompleted)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengubah status aktivitas"})
	}

	return c.JSON(fiber.Map{
		"id":           id,
		"is_completed": isCompleted,
	})
}
