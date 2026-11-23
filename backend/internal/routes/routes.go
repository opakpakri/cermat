package routes

import (
	"cermat-api/internal/handlers"
	"cermat-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(app *fiber.App, db *pgxpool.Pool, jwtSecret string) {
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "app": "Cermat API"})
	})

	// Auth routes (public)
	authHandler := &handlers.AuthHandler{DB: db, JWTSecret: jwtSecret}
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes (JWT + Tenant Context Middleware)
	protected := api.Group("", middleware.AuthTenant(jwtSecret))

	dashboardHandler := &handlers.DashboardHandler{DB: db}
	protected.Get("/dashboard/summary", dashboardHandler.GetSummary)

	dealHandler := &handlers.DealHandler{DB: db}
	protected.Get("/deals", dealHandler.List)
	protected.Post("/deals", dealHandler.Create)
	protected.Get("/deals/:id", dealHandler.GetByID)
	protected.Patch("/deals/:id/stage", dealHandler.UpdateStage)
	protected.Put("/deals/:id", dealHandler.Update)

	contactHandler := &handlers.ContactHandler{DB: db}
	protected.Get("/contacts", contactHandler.List)
	protected.Post("/contacts", contactHandler.Create)
	protected.Get("/contacts/:id", contactHandler.GetByID)
	protected.Put("/contacts/:id", contactHandler.Update)
	protected.Delete("/contacts/:id", contactHandler.Delete)

	catalogHandler := &handlers.CatalogHandler{DB: db}
	protected.Get("/catalog", catalogHandler.List)
	protected.Post("/catalog", catalogHandler.Create)
	protected.Put("/catalog/:id", catalogHandler.Update)
	protected.Patch("/catalog/:id/status", catalogHandler.UpdateStatus)
	protected.Delete("/catalog/:id", catalogHandler.Delete)

	activityHandler := &handlers.ActivityHandler{DB: db}
	protected.Get("/activities", activityHandler.List)
	protected.Post("/activities", activityHandler.Create)
	protected.Patch("/activities/:id/toggle", activityHandler.Toggle)

	ticketHandler := &handlers.TicketHandler{DB: db}
	protected.Get("/tickets", ticketHandler.List)
	protected.Post("/tickets", ticketHandler.Create)
	protected.Patch("/tickets/:id/status", ticketHandler.UpdateStatus)

	userHandler := &handlers.UserHandler{DB: db}
	protected.Get("/users/profile", userHandler.GetProfile)
	protected.Put("/users/profile", userHandler.UpdateProfile)

	// Admin-only User Management
	admin := protected.Group("/admin", middleware.RequireAdmin(db))
	admin.Get("/users", userHandler.AdminListUsers)
	admin.Post("/users", userHandler.AdminCreateUser)
	admin.Put("/users/:id", userHandler.AdminUpdateUser)
	admin.Delete("/users/:id", userHandler.AdminDeleteUser)
}
