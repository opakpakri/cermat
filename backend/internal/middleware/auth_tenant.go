package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthTenant(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token autentikasi tidak ditemukan",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Klaim token tidak valid",
			})
		}

		rawUserID, exists := claims["user_id"]
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Identitas tenant tidak ditemukan dalam token",
			})
		}

		var userID int
		switch v := rawUserID.(type) {
		case float64:
			userID = int(v)
		case int:
			userID = v
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Format identitas tenant tidak valid",
			})
		}

		if userID <= 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Identitas tenant tidak valid",
			})
		}

		// Inject verified userID strictly from JWT claims
		c.Locals("userID", userID)
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) int {
	if uid, ok := c.Locals("userID").(int); ok {
		return uid
	}
	return 0
}

func RequireAdmin(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := GetUserID(c)
		var role string
		err := db.QueryRow(context.Background(), `SELECT role FROM users WHERE id = $1`, userID).Scan(&role)
		if err != nil || role != "ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya Admin yang dapat mengakses fitur ini.",
			})
		}
		return c.Next()
	}
}
