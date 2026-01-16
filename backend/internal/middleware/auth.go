package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/service"
)

func AuthMiddleware(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{
				Error: domain.ErrorDetail{
					Code:    "UNAUTHORIZED",
					Message: "Token tidak ditemukan",
				},
			})
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{
				Error: domain.ErrorDetail{
					Code:    "UNAUTHORIZED",
					Message: "Format token tidak valid",
				},
			})
		}

		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{
				Error: domain.ErrorDetail{
					Code:    "UNAUTHORIZED",
					Message: "Token tidak valid atau sudah expired",
				},
			})
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}

func RoleMiddleware(roles ...domain.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*service.JWTClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{
				Error: domain.ErrorDetail{
					Code:    "UNAUTHORIZED",
					Message: "Token tidak valid",
				},
			})
		}

		for _, role := range roles {
			if claims.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(domain.ErrorResponse{
			Error: domain.ErrorDetail{
				Code:    "FORBIDDEN",
				Message: "Anda tidak memiliki akses ke resource ini",
			},
		})
	}
}
