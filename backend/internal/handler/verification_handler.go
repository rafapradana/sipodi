package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/service"
)

type VerificationHandler struct {
	talentService *service.TalentService
}

func NewVerificationHandler(talentService *service.TalentService) *VerificationHandler {
	return &VerificationHandler{talentService: talentService}
}

func (h *VerificationHandler) ListPending(c *fiber.Ctx) error {
	params := h.parseListParams(c)
	if params.Filters["status"] == "" {
		params.Filters["status"] = "pending"
	}

	// Admin sekolah can only see talents from their school
	claims := GetClaims(c)
	if claims.Role == domain.RoleAdminSekolah && claims.SchoolID != nil {
		params.Filters["school_id"] = claims.SchoolID.String()
	}

	talents, total, err := h.talentService.List(c.Context(), params)
	if err != nil {
		return InternalError(c)
	}

	var resp []fiber.Map
	for _, talent := range talents {
		user, _ := h.talentService.GetUser(c.Context(), talent.UserID)
		detail, _, _ := h.talentService.GetDetail(c.Context(), &talent)

		item := fiber.Map{
			"id":          talent.ID,
			"talent_type": talent.TalentType,
			"status":      talent.Status,
			"detail":      detail,
			"created_at":  talent.CreatedAt,
		}

		if user != nil {
			item["user"] = fiber.Map{
				"id":        user.ID,
				"full_name": user.FullName,
				"photo_url": user.PhotoURL,
			}
		}

		resp = append(resp, item)
	}

	meta := domain.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalCount:  total,
		TotalPages:  (total + params.Limit - 1) / params.Limit,
	}

	return SuccessList(c, resp, meta)
}

func (h *VerificationHandler) Approve(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	claims := GetClaims(c)

	// Check if admin sekolah can verify this talent
	if claims.Role == domain.RoleAdminSekolah {
		talent, err := h.talentService.GetByID(c.Context(), id)
		if err != nil {
			if err == service.ErrTalentNotFound {
				return NotFound(c, "Talenta tidak ditemukan")
			}
			return InternalError(c)
		}

		user, _ := h.talentService.GetUser(c.Context(), talent.UserID)
		if user == nil || user.SchoolID == nil || *user.SchoolID != *claims.SchoolID {
			return Forbidden(c, "Anda hanya dapat memverifikasi talenta GTK di sekolah Anda")
		}
	}

	talent, err := h.talentService.Approve(c.Context(), id, claims.UserID)
	if err != nil {
		switch err {
		case service.ErrTalentNotFound:
			return NotFound(c, "Talenta tidak ditemukan")
		case service.ErrAlreadyVerified:
			return BadRequest(c, "ALREADY_VERIFIED", "Talenta sudah diverifikasi sebelumnya")
		default:
			return InternalError(c)
		}
	}

	return SuccessWithMessage(c, fiber.Map{
		"id":          talent.ID,
		"status":      talent.Status,
		"verified_at": talent.VerifiedAt,
	}, "Talenta berhasil disetujui")
}

func (h *VerificationHandler) Reject(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	var req domain.RejectTalentRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	if req.RejectionReason == "" {
		return ValidationError(c, []domain.FieldError{
			{Field: "rejection_reason", Message: "Alasan penolakan wajib diisi"},
		})
	}

	claims := GetClaims(c)

	// Check if admin sekolah can verify this talent
	if claims.Role == domain.RoleAdminSekolah {
		talent, err := h.talentService.GetByID(c.Context(), id)
		if err != nil {
			if err == service.ErrTalentNotFound {
				return NotFound(c, "Talenta tidak ditemukan")
			}
			return InternalError(c)
		}

		user, _ := h.talentService.GetUser(c.Context(), talent.UserID)
		if user == nil || user.SchoolID == nil || *user.SchoolID != *claims.SchoolID {
			return Forbidden(c, "Anda hanya dapat memverifikasi talenta GTK di sekolah Anda")
		}
	}

	talent, err := h.talentService.Reject(c.Context(), id, claims.UserID, req.RejectionReason)
	if err != nil {
		switch err {
		case service.ErrTalentNotFound:
			return NotFound(c, "Talenta tidak ditemukan")
		case service.ErrAlreadyVerified:
			return BadRequest(c, "ALREADY_VERIFIED", "Talenta sudah diverifikasi sebelumnya")
		default:
			return InternalError(c)
		}
	}

	return SuccessWithMessage(c, fiber.Map{
		"id":               talent.ID,
		"status":           talent.Status,
		"rejection_reason": talent.RejectionReason,
		"verified_at":      talent.VerifiedAt,
	}, "Talenta ditolak")
}

func (h *VerificationHandler) BatchApprove(c *fiber.Ctx) error {
	var req domain.BatchApproveRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	claims := GetClaims(c)
	result, err := h.talentService.BatchApprove(c.Context(), req.IDs, claims.UserID)
	if err != nil {
		return InternalError(c)
	}

	message := ""
	if result.FailedCount > 0 {
		message = strconv.Itoa(result.SuccessCount) + " talenta berhasil disetujui, " + strconv.Itoa(result.FailedCount) + " gagal"
	} else {
		message = strconv.Itoa(result.SuccessCount) + " talenta berhasil disetujui"
	}

	return SuccessWithMessage(c, fiber.Map{
		"approved_count": result.SuccessCount,
		"failed_count":   result.FailedCount,
		"failed_ids":     result.FailedIDs,
	}, message)
}

func (h *VerificationHandler) BatchReject(c *fiber.Ctx) error {
	var req domain.BatchRejectRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	if req.RejectionReason == "" {
		return ValidationError(c, []domain.FieldError{
			{Field: "rejection_reason", Message: "Alasan penolakan wajib diisi"},
		})
	}

	claims := GetClaims(c)
	result, err := h.talentService.BatchReject(c.Context(), req.IDs, claims.UserID, req.RejectionReason)
	if err != nil {
		return InternalError(c)
	}

	message := ""
	if result.FailedCount > 0 {
		message = strconv.Itoa(result.SuccessCount) + " talenta ditolak, " + strconv.Itoa(result.FailedCount) + " gagal"
	} else {
		message = strconv.Itoa(result.SuccessCount) + " talenta ditolak"
	}

	return SuccessWithMessage(c, fiber.Map{
		"rejected_count": result.SuccessCount,
		"failed_count":   result.FailedCount,
		"failed_ids":     result.FailedIDs,
	}, message)
}

func (h *VerificationHandler) parseListParams(c *fiber.Ctx) domain.ListParams {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	return domain.ListParams{
		Page:  page,
		Limit: limit,
		Sort:  c.Query("sort"),
		Filters: map[string]string{
			"status":      c.Query("status"),
			"school_id":   c.Query("school_id"),
			"talent_type": c.Query("talent_type"),
		},
	}
}
