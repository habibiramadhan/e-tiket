//internal/delivery/http/routes/user_routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"
	
	"ticket-system/internal/delivery/http/handler"
	"ticket-system/internal/delivery/http/middleware"
)

func SetupUserRoutes(
	router fiber.Router,
	userHandler *handler.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
	loggerMiddleware *middleware.LoggerMiddleware,
) {
	// Public routes
	router.Post("/register", userHandler.Register)
	router.Post("/login", userHandler.Login)
	router.Get("/verify-email", userHandler.VerifyEmail)
	router.Post("/resend-verification", userHandler.ResendVerificationEmail)
	
	// Protected routes
	router.Put("/profile", authMiddleware.AuthenticateJWT(), userHandler.UpdateProfile)
}