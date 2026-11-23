package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"cermat-api/internal/config"
	"cermat-api/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	db, err := config.ConnectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// Run schema migration automatically if file exists
	migrationPath := filepath.Join("migrations", "001_init_multi_tenant.sql")
	if _, err := os.Stat(migrationPath); err == nil {
		if err := config.RunMigrations(ctx, db, migrationPath); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Database migration verified")
	}

	if err := config.EnsureAdminUser(ctx, db); err != nil {
		log.Fatalf("Admin seed failed: %v", err)
	}
	log.Println("Admin user verified (admin@cermat.id)")

	app := fiber.New(fiber.Config{
		AppName: "Cermat (Cellix Relationship Management Terpadu) API",
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	routes.Setup(app, db, cfg.JWTSecret)

	log.Printf("Cermat API listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
