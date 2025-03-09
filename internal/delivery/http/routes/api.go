//internal/delivery/http/routes/api.go (modified)
package routes

import (
	"database/sql"
	"fmt"
	"log"
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
	app.Use(recover.New())
	
	userRepo := postgres.NewUserRepository(db)
	userProfileRepo := postgres.NewUserProfileRepository(db)
	emailVerificationRepo := postgres.NewEmailVerificationRepository(db)
	eventRepo := postgres.NewEventRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	loggerMiddleware := middleware.NewLoggerMiddleware()
	
	smtpConfig := utils.SMTPConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		FromName: cfg.SMTPFromName,
	}
	
	appURL := fmt.Sprintf("http://localhost:%s", cfg.ServerPort)

	if cfg.AppEnv != "development" {
		appURL = "https://ea6f-111-94-16-128.ngrok-free.app/" 
	}
	
	userUsecase := usecase.NewUserUsecase(
		userRepo, 
		userProfileRepo, 
		emailVerificationRepo,
		cfg.JWTSecret, 
		cfg.TokenExpiry, 
		smtpConfig,
		appURL,
	)
	
	eventUsecase := usecase.NewEventUsecase(eventRepo, userRepo)
	
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, eventRepo, userRepo)
	
	userHandler := handler.NewUserHandler(userUsecase)
	eventHandler := handler.NewEventHandler(eventUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)
	
	api := app.Group("/api", loggerMiddleware.LogRequest())

	SetupUserRoutes(api, userHandler, authMiddleware, loggerMiddleware)
	SetupEventRoutes(api, eventHandler, authMiddleware)
	SetupTransactionRoutes(api, transactionHandler, authMiddleware)
	
	log.Println("Registered routes:")
	for _, r := range app.GetRoutes() {
		log.Printf("METHOD: %s, PATH: %s", r.Method, r.Path)
	}
	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Aplikasi Event Management API berjalan")
	})
	
	app.Use(func(c *fiber.Ctx) error {
		return utils.ErrorResponse(c, "Not Found", "Endpoint tidak ditemukan", fiber.StatusNotFound)
	})
}