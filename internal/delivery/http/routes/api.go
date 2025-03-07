//internal/delivery/http/routes/api.go

package routes

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	
	"ticket-system/internal/delivery/http/handler"
	"ticket-system/internal/delivery/http/middleware"
	"ticket-system/internal/repository/postgres"
	"ticket-system/internal/usecase"
	"ticket-system/pkg/config"
	"ticket-system/pkg/utils"
)

func SetupRoutes(app *fiber.App, db *sql.DB, cfg *config.Config) {
	// Middleware global
	app.Use(recover.New())
	
	// Repositories
	userRepo := postgres.NewUserRepository(db)
	userProfileRepo := postgres.NewUserProfileRepository(db)
	emailVerificationRepo := postgres.NewEmailVerificationRepository(db)
	
	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	loggerMiddleware := middleware.NewLoggerMiddleware()
	
	// SMTP Config
	smtpConfig := utils.SMTPConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		FromName: cfg.SMTPFromName,
	}
	
	// App URL untuk link verifikasi
	appURL := fmt.Sprintf("http://localhost:%s", cfg.ServerPort)
	if cfg.AppEnv != "development" {
		// Di production, gunakan domain yang sebenarnya
		appURL = "https://your-domain.com" // Sesuaikan dengan domain Anda
	}
	
	// Usecases
	userUsecase := usecase.NewUserUsecase(
		userRepo, 
		userProfileRepo, 
		emailVerificationRepo,
		cfg.JWTSecret, 
		cfg.TokenExpiry, 
		smtpConfig,
		appURL,
	)
	
	// Handlers
	userHandler := handler.NewUserHandler(userUsecase)
	
	// API Group dengan logger middleware
	api := app.Group("/api", loggerMiddleware.LogRequest())
	
	// Setup routes
	SetupUserRoutes(api, userHandler, authMiddleware, loggerMiddleware)
	
	// Default route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Aplikasi Event Management API berjalan")
	})
	
	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return utils.ErrorResponse(c, "Not Found", "Endpoint tidak ditemukan", fiber.StatusNotFound)
	})
}