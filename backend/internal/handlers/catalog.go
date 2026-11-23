package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"cermat-api/internal/middleware"
	"cermat-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogHandler struct {
	DB *pgxpool.Pool
}

type catalogRequest struct {
	Code     *string         `json:"code"`
	Name     string          `json:"name"`
	Category *string         `json:"category"`
	Price    float64         `json:"price"`
	Status   string          `json:"status"`
	Specs    json.RawMessage `json:"specs"`
}

func (h *CatalogHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	status := strings.ToUpper(strings.TrimSpace(c.Query("status")))
	search := strings.TrimSpace(c.Query("search"))

	query := `
		SELECT id, user_id, code, name, category, price, status, specs, created_at
		FROM catalog_items
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
		query += fmt.Sprintf(" AND (name ILIKE $%d OR COALESCE(code, '') ILIKE $%d OR COALESCE(category, '') ILIKE $%d)", idx, idx, idx)
		args = append(args, "%"+search+"%")
		idx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := h.DB.Query(c.Context(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat katalog"})
	}
	defer rows.Close()

	items := make([]models.CatalogItem, 0)
	for rows.Next() {
		var item models.CatalogItem
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.Code, &item.Name, &item.Category,
			&item.Price, &item.Status, &item.Specs, &item.CreatedAt,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca item katalog"})
		}
		items = append(items, item)
	}

	return c.JSON(items)
}

func (h *CatalogHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req catalogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))
	if req.Name == "" || req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama item dan harga valid wajib diisi"})
	}
	if req.Status == "" {
		req.Status = "AVAILABLE"
	}
	if len(req.Specs) == 0 {
		req.Specs = json.RawMessage(`{}`)
	}

	var item models.CatalogItem
	err := h.DB.QueryRow(c.Context(), `
		INSERT INTO catalog_items (user_id, code, name, category, price, status, specs)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, code, name, category, price, status, specs, created_at
	`, userID, req.Code, req.Name, req.Category, req.Price, req.Status, req.Specs).Scan(
		&item.ID, &item.UserID, &item.Code, &item.Name, &item.Category,
		&item.Price, &item.Status, &item.Specs, &item.CreatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menambah item katalog"})
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *CatalogHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item katalog tidak ditemukan"})
	}

	var req catalogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))
	if req.Name == "" || req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama item dan harga valid wajib diisi"})
	}
	if req.Status == "" {
		req.Status = "AVAILABLE"
	}
	if len(req.Specs) == 0 {
		req.Specs = json.RawMessage(`{}`)
	}

	var item models.CatalogItem
	err = h.DB.QueryRow(c.Context(), `
		UPDATE catalog_items
		SET code = $1, name = $2, category = $3, price = $4, status = $5, specs = $6
		WHERE id = $7 AND user_id = $8
		RETURNING id, user_id, code, name, category, price, status, specs, created_at
	`, req.Code, req.Name, req.Category, req.Price, req.Status, req.Specs, id, userID).Scan(
		&item.ID, &item.UserID, &item.Code, &item.Name, &item.Category,
		&item.Price, &item.Status, &item.Specs, &item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item katalog tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui item katalog"})
	}

	return c.JSON(item)
}

func (h *CatalogHandler) UpdateStatus(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item katalog tidak ditemukan"})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}
	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))
	if req.Status != "AVAILABLE" && req.Status != "BOOKED" && req.Status != "SOLD" && req.Status != "INACTIVE" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status tidak valid"})
	}

	var item models.CatalogItem
	err = h.DB.QueryRow(c.Context(), `
		UPDATE catalog_items
		SET status = $1
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, code, name, category, price, status, specs, created_at
	`, req.Status, id, userID).Scan(
		&item.ID, &item.UserID, &item.Code, &item.Name, &item.Category,
		&item.Price, &item.Status, &item.Specs, &item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item katalog tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui status item"})
	}

	return c.JSON(item)
}

func (h *CatalogHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item katalog tidak ditemukan"})
	}

	cmd, err := h.DB.Exec(c.Context(), `DELETE FROM catalog_items WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus item katalog"})
	}
	if cmd.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item katalog tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Item katalog berhasil dihapus"})
}
