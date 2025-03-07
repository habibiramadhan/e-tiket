//internal/delivery/http/middleware/error_middleware.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"ticket-system/pkg/utils"
)

type ErrorMiddleware struct{}

func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

func (m *ErrorMiddleware) ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		
		switch code {
		case fiber.StatusNotFound:
			return utils.NotFoundError(c, "Endpoint tidak ditemukan")
		case fiber.StatusUnauthorized:
			return utils.UnauthorizedError(c, "Anda tidak terautentikasi")
		case fiber.StatusForbidden:
			return utils.ErrorResponse(c, utils.ErrorCodeUnauthorized, "Anda tidak memiliki izin", fiber.StatusForbidden)
		case fiber.StatusMethodNotAllowed:
			return utils.ErrorResponse(c, "SRV005", "Metode HTTP tidak diizinkan", code)
		case fiber.StatusTooManyRequests:
			return utils.ErrorResponse(c, "SRV006", "Terlalu banyak permintaan", code)
		case fiber.StatusBadGateway, fiber.StatusGatewayTimeout, fiber.StatusServiceUnavailable:
			return utils.ErrorResponse(c, utils.ErrorCodeExternalServiceError, "Layanan tidak tersedia", code)
		default:
			return utils.ServerError(c, "Terjadi kesalahan pada server")
		}
	}
}