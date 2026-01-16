package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/service"
)

type TalentHandler struct {
	talentService *service.TalentService
	uploadService *service.UploadService
}

func NewTalentHandler(talentService *service.TalentService, uploadService *service.UploadService) *TalentHandler {
	return &TalentHandler{
		talentService: talentService,
		uploadService: uploadService,
	}
}

func (h *TalentHandler) List(c *fiber.Ctx) error {
	params := h.parseListParams(c)

	// Admin sekolah can only see talents from their school
	claims := GetClaims(c)
	if claims.Role == domain.RoleAdminSekolah && claims.SchoolID != nil {
		params.Filters["school_id"] = claims.SchoolID.String()
	}

	talents, total, err := h.talentService.List(c.Context(), params)
	if err != nil {
		return InternalError(c)
	}

	var resp []domain.TalentListResponse
	for _, talent := range talents {
		resp = append(resp, h.toTalentListResponse(c, &talent))
	}

	meta := domain.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalCount:  total,
		TotalPages:  (total + params.Limit - 1) / params.Limit,
	}

	return SuccessList(c, resp, meta)
}

func (h *TalentHandler) ListMyTalents(c *fiber.Ctx) error {
	claims := GetClaims(c)
	params := h.parseListParams(c)
	params.Filters["user_id"] = claims.UserID.String()

	talents, total, err := h.talentService.List(c.Context(), params)
	if err != nil {
		return InternalError(c)
	}

	var resp []domain.TalentListResponse
	for _, talent := range talents {
		resp = append(resp, h.toTalentListResponse(c, &talent))
	}

	meta := domain.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalCount:  total,
		TotalPages:  (total + params.Limit - 1) / params.Limit,
	}

	return SuccessList(c, resp, meta)
}

func (h *TalentHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	talent, err := h.talentService.GetByID(c.Context(), id)
	if err != nil {
		if err == service.ErrTalentNotFound {
			return NotFound(c, "Talenta tidak ditemukan")
		}
		return InternalError(c)
	}

	resp := h.toTalentResponse(c, talent)
	return Success(c, resp)
}

func (h *TalentHandler) Create(c *fiber.Ctx) error {
	claims := GetClaims(c)
	var req domain.CreateTalentRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	// Validate detail based on talent type
	if err := h.validateTalentDetail(req.TalentType, req.Detail); err != nil {
		return ValidationError(c, err)
	}

	// Get certificate URL if upload_id provided
	var certURL *string
	if req.UploadID != nil {
		info, ok := h.uploadService.GetUploadInfo(*req.UploadID)
		if ok {
			url := h.uploadService.GetFileURL(info.ObjectName)
			certURL = &url
		}
	}

	talent, err := h.talentService.Create(c.Context(), claims.UserID, req, certURL)
	if err != nil {
		return InternalError(c)
	}

	resp := h.toTalentResponse(c, talent)
	return SuccessCreated(c, resp, "Talenta berhasil ditambahkan dan menunggu verifikasi")
}

func (h *TalentHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	claims := GetClaims(c)
	var req domain.UpdateTalentRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	// Get certificate URL if upload_id provided
	var certURL *string
	if req.UploadID != nil {
		info, ok := h.uploadService.GetUploadInfo(*req.UploadID)
		if ok {
			url := h.uploadService.GetFileURL(info.ObjectName)
			certURL = &url
		}
	}

	talent, err := h.talentService.Update(c.Context(), id, claims.UserID, req, certURL)
	if err != nil {
		switch err {
		case service.ErrTalentNotFound:
			return NotFound(c, "Talenta tidak ditemukan")
		case service.ErrForbidden:
			return Forbidden(c, "Anda hanya dapat mengubah talenta milik sendiri")
		default:
			return InternalError(c)
		}
	}

	resp := h.toTalentResponse(c, talent)
	return SuccessWithMessage(c, resp, "Talenta berhasil diperbarui dan menunggu verifikasi ulang")
}

func (h *TalentHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	claims := GetClaims(c)
	err = h.talentService.Delete(c.Context(), id, claims.UserID)
	if err != nil {
		switch err {
		case service.ErrTalentNotFound:
			return NotFound(c, "Talenta tidak ditemukan")
		case service.ErrForbidden:
			return Forbidden(c, "Anda hanya dapat menghapus talenta milik sendiri")
		default:
			return InternalError(c)
		}
	}

	return SuccessNoContent(c)
}

func (h *TalentHandler) parseListParams(c *fiber.Ctx) domain.ListParams {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	return domain.ListParams{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
		Sort:   c.Query("sort"),
		Filters: map[string]string{
			"user_id":     c.Query("user_id"),
			"school_id":   c.Query("school_id"),
			"talent_type": c.Query("talent_type"),
			"status":      c.Query("status"),
		},
	}
}

func (h *TalentHandler) toTalentResponse(c *fiber.Ctx, talent *domain.Talent) domain.TalentResponse {
	resp := domain.TalentResponse{
		ID:              talent.ID,
		TalentType:      talent.TalentType,
		Status:          talent.Status,
		VerifiedAt:      talent.VerifiedAt,
		RejectionReason: talent.RejectionReason,
		CreatedAt:       talent.CreatedAt,
		UpdatedAt:       talent.UpdatedAt,
	}

	// Get user info
	user, _ := h.talentService.GetUser(c.Context(), talent.UserID)
	if user != nil {
		resp.User = &domain.UserRef{
			ID:       user.ID,
			FullName: user.FullName,
		}
	}

	// Get detail and certificate URL
	detail, certURL, _ := h.talentService.GetDetail(c.Context(), talent)
	resp.Detail = detail
	resp.CertificateURL = certURL

	// Get verifier info
	if talent.VerifiedBy != nil {
		verifier, _ := h.talentService.GetUser(c.Context(), *talent.VerifiedBy)
		if verifier != nil {
			resp.VerifiedBy = &domain.UserRef{
				ID:       verifier.ID,
				FullName: verifier.FullName,
			}
		}
	}

	return resp
}

func (h *TalentHandler) toTalentListResponse(c *fiber.Ctx, talent *domain.Talent) domain.TalentListResponse {
	resp := domain.TalentListResponse{
		ID:         talent.ID,
		TalentType: talent.TalentType,
		Status:     talent.Status,
		CreatedAt:  talent.CreatedAt,
		UpdatedAt:  talent.UpdatedAt,
	}

	// Get user info
	user, _ := h.talentService.GetUser(c.Context(), talent.UserID)
	if user != nil {
		schoolName := ""
		if user.SchoolID != nil {
			// Would need school service here, simplified for now
		}
		resp.User = &domain.TalentUser{
			ID:         user.ID,
			FullName:   user.FullName,
			SchoolName: schoolName,
		}
	}

	// Get detail
	detail, _, _ := h.talentService.GetDetail(c.Context(), talent)
	resp.Detail = detail

	return resp
}

func (h *TalentHandler) validateTalentDetail(talentType domain.TalentType, detail interface{}) []domain.FieldError {
	var errors []domain.FieldError

	if detail == nil {
		errors = append(errors, domain.FieldError{Field: "detail", Message: "Detail wajib diisi"})
		return errors
	}

	// Convert detail to map for validation
	detailMap, ok := detail.(map[string]interface{})
	if !ok {
		errors = append(errors, domain.FieldError{Field: "detail", Message: "Format detail tidak valid"})
		return errors
	}

	switch talentType {
	case domain.TalentTypePesertaPelatihan:
		if detailMap["activity_name"] == nil || detailMap["activity_name"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.activity_name", Message: "Nama kegiatan wajib diisi"})
		}
		if detailMap["organizer"] == nil || detailMap["organizer"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.organizer", Message: "Penyelenggara wajib diisi"})
		}
		if detailMap["start_date"] == nil || detailMap["start_date"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.start_date", Message: "Tanggal mulai wajib diisi"})
		}
		if detailMap["duration_days"] == nil {
			errors = append(errors, domain.FieldError{Field: "detail.duration_days", Message: "Jangka waktu wajib diisi"})
		}

	case domain.TalentTypePembimbingLomba:
		if detailMap["competition_name"] == nil || detailMap["competition_name"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.competition_name", Message: "Nama lomba wajib diisi"})
		}
		if detailMap["level"] == nil || detailMap["level"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.level", Message: "Jenjang wajib diisi"})
		}
		if detailMap["organizer"] == nil || detailMap["organizer"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.organizer", Message: "Penyelenggara wajib diisi"})
		}
		if detailMap["field"] == nil || detailMap["field"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.field", Message: "Bidang wajib diisi"})
		}
		if detailMap["achievement"] == nil || detailMap["achievement"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.achievement", Message: "Prestasi wajib diisi"})
		}

	case domain.TalentTypePesertaLomba:
		if detailMap["competition_name"] == nil || detailMap["competition_name"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.competition_name", Message: "Nama lomba wajib diisi"})
		}
		if detailMap["level"] == nil || detailMap["level"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.level", Message: "Jenjang wajib diisi"})
		}
		if detailMap["organizer"] == nil || detailMap["organizer"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.organizer", Message: "Penyelenggara wajib diisi"})
		}
		if detailMap["field"] == nil || detailMap["field"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.field", Message: "Bidang wajib diisi"})
		}
		if detailMap["start_date"] == nil || detailMap["start_date"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.start_date", Message: "Tanggal mulai wajib diisi"})
		}
		if detailMap["competition_field"] == nil || detailMap["competition_field"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.competition_field", Message: "Bidang lomba wajib diisi"})
		}
		if detailMap["achievement"] == nil || detailMap["achievement"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.achievement", Message: "Prestasi wajib diisi"})
		}

	case domain.TalentTypeMinatBakat:
		if detailMap["interest_name"] == nil || detailMap["interest_name"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.interest_name", Message: "Nama minat/bakat wajib diisi"})
		}
		if detailMap["description"] == nil || detailMap["description"] == "" {
			errors = append(errors, domain.FieldError{Field: "detail.description", Message: "Deskripsi wajib diisi"})
		}
	}

	return errors
}
