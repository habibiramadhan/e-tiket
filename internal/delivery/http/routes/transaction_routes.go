//internal/delivery/http/routes/transaction_routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"
	
	"ticket-system/internal/delivery/http/handler"
	"ticket-system/internal/delivery/http/middleware"
)

func SetupTransactionRoutes(
	router fiber.Router,
	transactionHandler *handler.TransactionHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	transactionRoutes := router.Group("/transactions")
	transactionRoutes.Use(authMiddleware.AuthenticateJWT())

	transactionRoutes.Get("/code", transactionHandler.GetTransactionByCode)
	transactionRoutes.Post("/proof", transactionHandler.UploadPaymentProof)
	transactionRoutes.Put("/:id/cancel", transactionHandler.CancelTransaction)
	transactionRoutes.Get("/:id", transactionHandler.GetTransactionByID)
	transactionRoutes.Get("", transactionHandler.GetUserTransactions)
	transactionRoutes.Post("", transactionHandler.CreateTransaction)

	organizerRoutes := router.Group("/organizer/transactions")
	organizerRoutes.Use(authMiddleware.AuthenticateJWT())
	organizerRoutes.Use(authMiddleware.RoleCheck([]string{"organizer"}))

	organizerRoutes.Put("/:id/verify", transactionHandler.VerifyPayment)
}