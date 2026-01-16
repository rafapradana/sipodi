package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/service"
)

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

func (h *UploadHandler) Presign(c *fiber.Ctx) error {
	claims := GetClaims(c)
	var req domain.PresignRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	// Validation
	var errors []domain.FieldError
	if req.Filename == "" {
		errors = append(errors, domain.FieldError{Field: "filename", Message: "Nama file wajib diisi"})
	}
	if req.ContentType == "" {
		errors = append(errors, domain.FieldError{Field: "content_type", Message: "Content type wajib diisi"})
	}
	if req.UploadType == "" {
		errors = append(errors, domain.FieldError{Field: "upload_type", Message: "Upload type wajib diisi"})
	}
	if len(errors) > 0 {
		return ValidationError(c, errors)
	}

	resp, err := h.uploadService.GeneratePresignedURL(c.Context(), claims.UserID, req)
	if err != nil {
		if err == service.ErrInvalidFileType {
			return BadRequest(c, "INVALID_FILE_TYPE", "Tipe file tidak diizinkan. Gunakan PDF atau gambar (JPG, PNG)")
		}
		return InternalError(c)
	}

	return Success(c, resp)
}

func (h *UploadHandler) Confirm(c *fiber.Ctx) error {
	uploadID, err := uuid.Parse(c.Params("upload_id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "Upload ID tidak valid")
	}

	claims := GetClaims(c)
	resp, err := h.uploadService.ConfirmUpload(c.Context(), uploadID, claims.UserID)
	if err != nil {
		switch err {
		case service.ErrUploadNotFound:
			return NotFound(c, "Upload tidak ditemukan atau sudah expired")
		case service.ErrFileNotUploaded:
			return BadRequest(c, "FILE_NOT_UPLOADED", "File belum diupload ke storage")
		default:
			return InternalError(c)
		}
	}

	return SuccessWithMessage(c, resp, "Upload berhasil dikonfirmasi")
}

func (h *UploadHandler) Cancel(c *fiber.Ctx) error {
	uploadID, err := uuid.Parse(c.Params("upload_id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "Upload ID tidak valid")
	}

	claims := GetClaims(c)
	err = h.uploadService.CancelUpload(c.Context(), uploadID, claims.UserID)
	if err != nil {
		if err == service.ErrUploadNotFound {
			return NotFound(c, "Upload tidak ditemukan")
		}
		return InternalError(c)
	}

	return SuccessNoContent(c)
}
