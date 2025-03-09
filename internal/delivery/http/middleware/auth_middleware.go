//internal/delivery/http/middleware/auth_middleware.go

package middleware

import (
	"log"
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
			log.Println("Auth failed: No Authorization header")
			return utils.ErrorResponse(c, utils.ErrorCodeTokenMissing, "Otorisasi diperlukan", fiber.StatusUnauthorized)
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			log.Println("Auth failed: Invalid token format")
			return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Format token tidak valid. Format: Bearer [token]", fiber.StatusUnauthorized)
		}

		tokenStr := splitToken[1]
		log.Printf("Validating token: %s", tokenStr[:10]+"...")
		
		claims, err := utils.ValidateToken(tokenStr, m.jwtSecret)
		if err != nil {
			log.Printf("Auth failed: Token validation error: %v", err)
			return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid: "+err.Error(), fiber.StatusUnauthorized)
		}

		log.Printf("Token valid, user: %s, role: %s", claims.Username, claims.Role)
		c.Locals("claims", claims)
		return c.Next()
	}
}

func (m *AuthMiddleware) RoleCheck(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*utils.JWTClaim)
		if !ok {
			log.Println("RoleCheck failed: No claims found in context")
			return utils.ErrorResponse(c, utils.ErrorCodeTokenMissing, "Token diperlukan", fiber.StatusUnauthorized)
		}

		userRole := claims.Role
		allowed := false
		log.Printf("Checking role: User has '%s', required one of: %v", userRole, roles)
		
		for _, role := range roles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			log.Printf("RoleCheck failed: User role '%s' not allowed (required: %v)", userRole, roles)
			return utils.ErrorResponse(c, utils.ErrorCodeUnauthorized, "Anda tidak memiliki izin untuk mengakses resource ini. Role Anda: "+userRole, fiber.StatusForbidden)
		}

		log.Println("RoleCheck passed")
		return c.Next()
	}
}