package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sipodi/backend/internal/domain"
)

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(domain.DataResponse{Data: data})
}

func SuccessWithMessage(c *fiber.Ctx, data interface{}, message string) error {
	return c.JSON(domain.DataResponse{Data: data, Message: message})
}

func SuccessCreated(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(domain.DataResponse{Data: data, Message: message})
}

func SuccessNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func SuccessList(c *fiber.Ctx, data interface{}, meta domain.PaginationMeta) error {
	return c.JSON(domain.ListResponse{Data: data, Meta: meta})
}

func Message(c *fiber.Ctx, message string) error {
	return c.JSON(fiber.Map{"message": message})
}

func Error(c *fiber.Ctx, status int, code string, message string) error {
	return c.Status(status).JSON(domain.ErrorResponse{
		Error: domain.ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func ErrorWithDetails(c *fiber.Ctx, status int, code string, message string, details []domain.FieldError) error {
	return c.Status(status).JSON(domain.ErrorResponse{
		Error: domain.ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func BadRequest(c *fiber.Ctx, code string, message string) error {
	return Error(c, fiber.StatusBadRequest, code, message)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, "NOT_FOUND", message)
}

func Conflict(c *fiber.Ctx, code string, message string) error {
	return Error(c, fiber.StatusConflict, code, message)
}

func ValidationError(c *fiber.Ctx, details []domain.FieldError) error {
	return ErrorWithDetails(c, fiber.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validasi gagal", details)
}

func InternalError(c *fiber.Ctx) error {
	return Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan pada server")
}
