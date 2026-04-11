package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/usecase"
	"github.com/l0ng7h0r/golang/pkg/security"
)

type AuthMiddleware struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthMiddleware(authUsecase *usecase.AuthUsecase) *AuthMiddleware {
	return &AuthMiddleware{authUsecase: authUsecase}
}

// Auth

func (m *AuthMiddleware) Auth(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := security.ParseToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "ParseToken failed: " + err.Error()})
	}

	c.Locals("user_id", (*claims)["user_id"])
	c.Locals("roles", (*claims)["roles"])

	return c.Next()
}

// RequireRole

func (m *AuthMiddleware) RequireRole(role string) fiber.Handler {
	return func(c fiber.Ctx) error {
		rolesLocals := c.Locals("roles")
		if rolesLocals == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "rolesLocals is nil (Auth middleware might not have set it)"})
		}

		var roles []string
		switch v := rolesLocals.(type) {
		case []string:
			roles = v
		case []interface{}:
			for _, r := range v {
				if strRole, ok := r.(string); ok {
					roles = append(roles, strRole)
				}
			}
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid roles format"})
		}

		if !contains(roles, role) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
		}
		return c.Next()
	}
}

// contains

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}