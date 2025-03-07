//internal/delivery/http/middleware/auth_middleware.go

package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"ticket-system/pkg/utils"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

func (m *AuthMiddleware) AuthenticateJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, utils.ErrorCodeTokenMissing, "Otorisasi diperlukan", fiber.StatusUnauthorized)
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Format token tidak valid. Format: Bearer [token]", fiber.StatusUnauthorized)
		}

		tokenStr := splitToken[1]
		claims, err := utils.ValidateToken(tokenStr, m.jwtSecret)
		if err != nil {
			return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid: "+err.Error(), fiber.StatusUnauthorized)
		}

		// Tambahkan claims ke locals (seperti context di net/http)
		c.Locals("claims", claims)
		return c.Next()
	}
}

func (m *AuthMiddleware) RoleCheck(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*utils.JWTClaim)
		if !ok {
			return utils.ErrorResponse(c, utils.ErrorCodeTokenMissing, "Token diperlukan", fiber.StatusUnauthorized)
		}

		userRole := claims.Role
		allowed := false
		for _, role := range roles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			return utils.ErrorResponse(c, utils.ErrorCodeUnauthorized, "Anda tidak memiliki izin untuk mengakses resource ini", fiber.StatusForbidden)
		}

		return c.Next()
	}
}