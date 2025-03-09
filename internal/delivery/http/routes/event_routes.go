//internal/delivery/http/routes/event_routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"
	
	"ticket-system/internal/delivery/http/handler"
	"ticket-system/internal/delivery/http/middleware"
)

func SetupEventRoutes(
	router fiber.Router,
	eventHandler *handler.EventHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Public routes 
	router.Get("/events", eventHandler.GetEventList)
	router.Get("/events/:id", eventHandler.GetEventByID)
	
	// Protected routes 
	organizerRoutes := router.Group("/organizer")
	organizerRoutes.Use(authMiddleware.AuthenticateJWT())
	organizerRoutes.Use(authMiddleware.RoleCheck([]string{"organizer"}))
	
	organizerRoutes.Post("/events", eventHandler.CreateEvent)
	organizerRoutes.Get("/events", eventHandler.GetEventsByOrganizer)
	organizerRoutes.Put("/events/:id", eventHandler.UpdateEvent)
	organizerRoutes.Delete("/events/:id", eventHandler.DeleteEvent)
	organizerRoutes.Get("/events/:id/sales", eventHandler.GetEventSales)
	
}