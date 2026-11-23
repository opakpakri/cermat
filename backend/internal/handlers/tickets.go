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

type TicketHandler struct {
	DB *pgxpool.Pool
}

type ticketRequest struct {
	ContactID   int     `json:"contact_id"`
	DealID      *int    `json:"deal_id"`
	Subject     string  `json:"subject"`
	Description *string `json:"description"`
	Priority    string  `json:"priority"`
}

func (h *TicketHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	status := strings.ToUpper(strings.TrimSpace(c.Query("status")))

	query := `
		SELECT t.id, t.user_id, t.ticket_code, t.contact_id, c.name, t.deal_id, d.title,
		       t.subject, t.description, t.priority, t.status, t.created_at, t.updated_at
		FROM tickets t
		JOIN contacts c ON c.id = t.contact_id AND c.user_id = t.user_id
		LEFT JOIN deals d ON d.id = t.deal_id AND d.user_id = t.user_id
		WHERE t.user_id = $1
	`
	args := []interface{}{userID}
	if status != "" {
		query += " AND t.status = $2"
		args = append(args, status)
	}
	query += " ORDER BY t.created_at DESC"

	rows, err := h.DB.Query(c.Context(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat tiket"})
	}
	defer rows.Close()

	tickets := make([]models.Ticket, 0)
	for rows.Next() {
		var tk models.Ticket
		if err := rows.Scan(
			&tk.ID, &tk.UserID, &tk.TicketCode, &tk.ContactID, &tk.ContactName, &tk.DealID, &tk.DealTitle,
			&tk.Subject, &tk.Description, &tk.Priority, &tk.Status, &tk.CreatedAt, &tk.UpdatedAt,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data tiket"})
		}
		tickets = append(tickets, tk)
	}

	return c.JSON(tickets)
}

func (h *TicketHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req ticketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Subject = strings.TrimSpace(req.Subject)
	req.Priority = strings.ToUpper(strings.TrimSpace(req.Priority))
	if req.Subject == "" || req.ContactID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Subjek tiket dan kontak wajib diisi"})
	}
	if req.Priority == "" {
		req.Priority = "MEDIUM"
	}

	// Cross-tenant validation: contact_id
	var contactExists bool
	err := h.DB.QueryRow(c.Context(), `SELECT EXISTS(SELECT 1 FROM contacts WHERE id = $1 AND user_id = $2)`, req.ContactID, userID).Scan(&contactExists)
	if err != nil || !contactExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kontak tidak ditemukan"})
	}

	// Cross-tenant validation: deal_id
	if req.DealID != nil && *req.DealID > 0 {
		var dealExists bool
		err = h.DB.QueryRow(c.Context(), `SELECT EXISTS(SELECT 1 FROM deals WHERE id = $1 AND user_id = $2)`, *req.DealID, userID).Scan(&dealExists)
		if err != nil || !dealExists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Deal tidak ditemukan"})
		}
	} else {
		req.DealID = nil
	}

	var count int
	_ = h.DB.QueryRow(c.Context(), `SELECT COUNT(*) FROM tickets WHERE user_id = $1`, userID).Scan(&count)
	ticketCode := fmt.Sprintf("TKT-%s-%03d", time.Now().Format("20060102"), count+1)

	var ticketID int
	err = h.DB.QueryRow(c.Context(), `
		INSERT INTO tickets (user_id, ticket_code, contact_id, deal_id, subject, description, priority)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, userID, ticketCode, req.ContactID, req.DealID, req.Subject, req.Description, req.Priority).Scan(&ticketID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat tiket"})
	}

	var tk models.Ticket
	err = h.DB.QueryRow(c.Context(), `
		SELECT t.id, t.user_id, t.ticket_code, t.contact_id, c.name, t.deal_id, d.title,
		       t.subject, t.description, t.priority, t.status, t.created_at, t.updated_at
		FROM tickets t
		JOIN contacts c ON c.id = t.contact_id AND c.user_id = t.user_id
		LEFT JOIN deals d ON d.id = t.deal_id AND d.user_id = t.user_id
		WHERE t.id = $1 AND t.user_id = $2
	`, ticketID, userID).Scan(
		&tk.ID, &tk.UserID, &tk.TicketCode, &tk.ContactID, &tk.ContactName, &tk.DealID, &tk.DealTitle,
		&tk.Subject, &tk.Description, &tk.Priority, &tk.Status, &tk.CreatedAt, &tk.UpdatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat tiket"})
	}

	return c.Status(fiber.StatusCreated).JSON(tk)
}

func (h *TicketHandler) UpdateStatus(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tiket tidak ditemukan"})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}
	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))
	if req.Status != "OPEN" && req.Status != "IN_PROGRESS" && req.Status != "RESOLVED" && req.Status != "CLOSED" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status tiket tidak valid"})
	}

	var currentStatus string
	err = h.DB.QueryRow(c.Context(), `
		SELECT status FROM tickets WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tiket tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memeriksa status tiket"})
	}

	if currentStatus == "CLOSED" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tiket yang sudah berstatus CLOSED tidak dapat diubah lagi",
		})
	}

	var tk models.Ticket
	err = h.DB.QueryRow(c.Context(), `
		UPDATE tickets
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND user_id = $3
		RETURNING id, status, updated_at
	`, req.Status, id, userID).Scan(&tk.ID, &tk.Status, &tk.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengubah status tiket"})
	}

	return c.JSON(fiber.Map{
		"id":         tk.ID,
		"status":     tk.Status,
		"updated_at": tk.UpdatedAt,
	})
}
