package handlers

import (
	"errors"
	"strconv"
	"strings"

	"cermat-api/internal/middleware"
	"cermat-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB *pgxpool.Pool
}

type updateProfileRequest struct {
	Name         string  `json:"name"`
	BusinessName string  `json:"business_name"`
	BusinessMode string  `json:"business_mode"`
	Phone        *string `json:"phone"`
	Address      *string `json:"address"`
	NewPassword  *string `json:"new_password"`
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var user models.User

	err := h.DB.QueryRow(c.Context(), `
		SELECT id, name, email, role, business_name, business_mode, phone, address, created_at, updated_at
		FROM users WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.BusinessMode,
		&user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Profil tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat profil"})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.BusinessName = strings.TrimSpace(req.BusinessName)
	req.BusinessMode = strings.ToUpper(strings.TrimSpace(req.BusinessMode))

	if req.Name == "" || req.BusinessName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama dan nama bisnis wajib diisi"})
	}
	if req.BusinessMode != "PROPERTY" && req.BusinessMode != "OFFICE" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Mode workspace harus PROPERTY atau OFFICE"})
	}

	var user models.User
	if req.NewPassword != nil && strings.TrimSpace(*req.NewPassword) != "" {
		if len(*req.NewPassword) < 6 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password baru minimal 6 karakter"})
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses password"})
		}
		err = h.DB.QueryRow(c.Context(), `
			UPDATE users
			SET name = $1, business_name = $2, business_mode = $3, phone = $4, address = $5,
			    password_hash = $6, updated_at = CURRENT_TIMESTAMP
			WHERE id = $7
			RETURNING id, name, email, role, business_name, business_mode, phone, address, created_at, updated_at
		`, req.Name, req.BusinessName, req.BusinessMode, req.Phone, req.Address, string(hash), userID).Scan(
			&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.BusinessMode,
			&user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui profil"})
		}
	} else {
		err := h.DB.QueryRow(c.Context(), `
			UPDATE users
			SET name = $1, business_name = $2, business_mode = $3, phone = $4, address = $5,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $6
			RETURNING id, name, email, role, business_name, business_mode, phone, address, created_at, updated_at
		`, req.Name, req.BusinessName, req.BusinessMode, req.Phone, req.Address, userID).Scan(
			&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.BusinessMode,
			&user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui profil"})
		}
	}

	return c.JSON(user)
}

func normalizeRole(r string) string {
	if strings.ToUpper(strings.TrimSpace(r)) == "ADMIN" {
		return "ADMIN"
	}
	return "USER"
}

func normalizeMode(m string) string {
	if strings.ToUpper(strings.TrimSpace(m)) == "PROPERTY" {
		return "PROPERTY"
	}
	return "OFFICE"
}

func (h *UserHandler) AdminListUsers(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	query := `
		SELECT id, name, email, role, business_name, business_mode, phone, address, created_at, updated_at
		FROM users
	`
	args := make([]interface{}, 0)
	if search != "" {
		query += ` WHERE name ILIKE $1 OR email ILIKE $1 OR business_name ILIKE $1`
		args = append(args, "%"+search+"%")
	}
	query += ` ORDER BY created_at DESC, id DESC`

	rows, err := h.DB.Query(c.Context(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar pengguna"})
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Role, &u.BusinessName, &u.BusinessMode,
			&u.Phone, &u.Address, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data pengguna"})
		}
		users = append(users, u)
	}

	return c.JSON(users)
}

func (h *UserHandler) AdminCreateUser(c *fiber.Ctx) error {
	var req models.AdminUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.BusinessName = strings.TrimSpace(req.BusinessName)
	req.Role = normalizeRole(req.Role)
	req.BusinessMode = normalizeMode(req.BusinessMode)

	if req.Name == "" || req.Email == "" || len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama, email, dan password (minimal 6 karakter) wajib diisi",
		})
	}
	if req.BusinessName == "" {
		req.BusinessName = "Cermat - Cellix CRM"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses password"})
	}

	var user models.User
	err = h.DB.QueryRow(c.Context(), `
		INSERT INTO users (name, email, password_hash, role, business_name, business_mode, phone, address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, email, role, business_name, business_mode, phone, address, created_at, updated_at
	`, req.Name, req.Email, string(hash), req.Role, req.BusinessName, req.BusinessMode, req.Phone, req.Address).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.BusinessMode,
		&user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email sudah terdaftar"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat pengguna baru"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) AdminUpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pengguna tidak valid"})
	}

	var req models.AdminUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.BusinessName = strings.TrimSpace(req.BusinessName)
	req.Role = normalizeRole(req.Role)
	req.BusinessMode = normalizeMode(req.BusinessMode)

	if req.Name == "" || req.Email == "" || req.BusinessName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama, email, dan nama bisnis wajib diisi"})
	}

	var user models.User
	if strings.TrimSpace(req.Password) != "" {
		if len(req.Password) < 6 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password baru minimal 6 karakter"})
		}
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses password"})
		}
		err = h.DB.QueryRow(c.Context(), `
			UPDATE users
			SET name = $1, email = $2, role = $3, business_name = $4, business_mode = $5,
			    phone = $6, address = $7, password_hash = $8, updated_at = CURRENT_TIMESTAMP
			WHERE id = $9
			RETURNING id, name, email, role, business_name, business_mode, phone, address, created_at, updated_at
		`, req.Name, req.Email, req.Role, req.BusinessName, req.BusinessMode, req.Phone, req.Address, string(hash), id).Scan(
			&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.BusinessMode,
			&user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		err = h.DB.QueryRow(c.Context(), `
			UPDATE users
			SET name = $1, email = $2, role = $3, business_name = $4, business_mode = $5,
			    phone = $6, address = $7, updated_at = CURRENT_TIMESTAMP
			WHERE id = $8
			RETURNING id, name, email, role, business_name, business_mode, phone, address, created_at, updated_at
		`, req.Name, req.Email, req.Role, req.BusinessName, req.BusinessMode, req.Phone, req.Address, id).Scan(
			&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.BusinessMode,
			&user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
		)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email sudah digunakan akun lain"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui pengguna"})
	}

	return c.JSON(user)
}

func (h *UserHandler) AdminDeleteUser(c *fiber.Ctx) error {
	currentUserID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pengguna tidak valid"})
	}
	if id == currentUserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tidak dapat menghapus akun Anda sendiri yang sedang aktif"})
	}

	tag, err := h.DB.Exec(c.Context(), `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus pengguna"})
	}
	if tag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Pengguna berhasil dihapus"})
}
