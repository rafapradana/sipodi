package handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	// Validation
	var errors []domain.FieldError
	if req.Email == "" {
		errors = append(errors, domain.FieldError{Field: "email", Message: "Email wajib diisi"})
	}
	if req.Password == "" {
		errors = append(errors, domain.FieldError{Field: "password", Message: "Password wajib diisi"})
	}
	if len(errors) > 0 {
		return ValidationError(c, errors)
	}

	resp, refreshToken, err := h.authService.Login(c.Context(), req)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			return Error(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "Email atau password salah")
		case service.ErrAccountDisabled:
			return Error(c, fiber.StatusForbidden, "ACCOUNT_DISABLED", "Akun Anda telah dinonaktifkan. Hubungi admin.")
		default:
			return InternalError(c)
		}
	}

	// Set refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	return Success(c, resp)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return Error(c, fiber.StatusUnauthorized, "INVALID_TOKEN", "Refresh token tidak ditemukan")
	}

	resp, newRefreshToken, err := h.authService.RefreshToken(c.Context(), refreshToken)
	if err != nil {
		switch err {
		case service.ErrTokenExpired:
			return Error(c, fiber.StatusUnauthorized, "TOKEN_EXPIRED", "Refresh token telah expired. Silakan login ulang.")
		case service.ErrInvalidToken:
			return Error(c, fiber.StatusUnauthorized, "INVALID_TOKEN", "Refresh token tidak valid")
		case service.ErrAccountDisabled:
			return Error(c, fiber.StatusForbidden, "ACCOUNT_DISABLED", "Akun Anda telah dinonaktifkan")
		default:
			return InternalError(c)
		}
	}

	// Set new refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		MaxAge:   7 * 24 * 60 * 60,
	})

	return Success(c, resp)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken != "" {
		h.authService.Logout(c.Context(), refreshToken)
	}

	// Clear cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-time.Hour),
	})

	return Message(c, "Berhasil logout")
}

func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	claims := GetClaims(c)
	count, err := h.authService.LogoutAll(c.Context(), claims.UserID)
	if err != nil {
		return InternalError(c)
	}

	// Clear cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-time.Hour),
	})

	return SuccessWithMessage(c, fiber.Map{"sessions_terminated": count}, "Berhasil logout dari semua perangkat")
}

// Helper to get claims from context
func GetClaims(c *fiber.Ctx) *service.JWTClaims {
	claims, _ := c.Locals("claims").(*service.JWTClaims)
	return claims
}

// Helper to extract token from header
func ExtractToken(c *fiber.Ctx) string {
	auth := c.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}
