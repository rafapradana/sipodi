package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/service"
)

type SchoolHandler struct {
	schoolService *service.SchoolService
}

func NewSchoolHandler(schoolService *service.SchoolService) *SchoolHandler {
	return &SchoolHandler{schoolService: schoolService}
}

func (h *SchoolHandler) List(c *fiber.Ctx) error {
	params := h.parseListParams(c)
	schools, total, err := h.schoolService.List(c.Context(), params)
	if err != nil {
		return InternalError(c)
	}

	var resp []domain.SchoolResponse
	for _, school := range schools {
		resp = append(resp, h.toSchoolResponse(c, &school))
	}

	meta := domain.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalCount:  total,
		TotalPages:  (total + params.Limit - 1) / params.Limit,
	}

	return SuccessList(c, resp, meta)
}

func (h *SchoolHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	school, err := h.schoolService.GetByID(c.Context(), id)
	if err != nil {
		if err == service.ErrSchoolNotFound {
			return NotFound(c, "Sekolah tidak ditemukan")
		}
		return InternalError(c)
	}

	resp := h.toSchoolDetailResponse(c, school)
	return Success(c, resp)
}

func (h *SchoolHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateSchoolRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	// Validation
	var errors []domain.FieldError
	if req.Name == "" {
		errors = append(errors, domain.FieldError{Field: "name", Message: "Nama sekolah wajib diisi"})
	}
	if req.NPSN == "" {
		errors = append(errors, domain.FieldError{Field: "npsn", Message: "NPSN wajib diisi"})
	}
	if req.Status != domain.SchoolStatusNegeri && req.Status != domain.SchoolStatusSwasta {
		errors = append(errors, domain.FieldError{Field: "status", Message: "Status harus negeri atau swasta"})
	}
	if req.Address == "" {
		errors = append(errors, domain.FieldError{Field: "address", Message: "Alamat wajib diisi"})
	}
	if len(errors) > 0 {
		return ValidationError(c, errors)
	}

	school, err := h.schoolService.Create(c.Context(), req)
	if err != nil {
		if err == service.ErrDuplicateNPSN {
			return Conflict(c, "DUPLICATE_NPSN", "NPSN sudah terdaftar")
		}
		return InternalError(c)
	}

	resp := h.toSchoolResponse(c, school)
	return SuccessCreated(c, resp, "Sekolah berhasil ditambahkan")
}

func (h *SchoolHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	var req domain.UpdateSchoolRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	school, err := h.schoolService.Update(c.Context(), id, req)
	if err != nil {
		switch err {
		case service.ErrSchoolNotFound:
			return NotFound(c, "Sekolah tidak ditemukan")
		case service.ErrDuplicateNPSN:
			return Conflict(c, "DUPLICATE_NPSN", "NPSN sudah terdaftar")
		case service.ErrInvalidHeadMaster:
			return BadRequest(c, "INVALID_HEAD_MASTER", "User yang dipilih bukan kepala sekolah")
		default:
			return InternalError(c)
		}
	}

	resp := h.toSchoolResponse(c, school)
	return SuccessWithMessage(c, resp, "Sekolah berhasil diperbarui")
}

func (h *SchoolHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	err = h.schoolService.Delete(c.Context(), id)
	if err != nil {
		switch err {
		case service.ErrSchoolNotFound:
			return NotFound(c, "Sekolah tidak ditemukan")
		case service.ErrSchoolHasUsers:
			return BadRequest(c, "HAS_DEPENDENCIES", "Tidak dapat menghapus sekolah yang masih memiliki GTK")
		default:
			return InternalError(c)
		}
	}

	return SuccessNoContent(c)
}

func (h *SchoolHandler) GetUsers(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	// Check authorization for admin_sekolah
	claims := GetClaims(c)
	if claims.Role == domain.RoleAdminSekolah && (claims.SchoolID == nil || *claims.SchoolID != id) {
		return Forbidden(c, "Anda hanya dapat melihat GTK di sekolah Anda")
	}

	params := h.parseListParams(c)
	users, total, err := h.schoolService.GetUsers(c.Context(), id, params)
	if err != nil {
		return InternalError(c)
	}

	var resp []fiber.Map
	for _, user := range users {
		resp = append(resp, fiber.Map{
			"id":        user.ID,
			"full_name": user.FullName,
			"nuptk":     user.NUPTK,
			"nip":       user.NIP,
			"gtk_type":  user.GTKType,
			"position":  user.Position,
			"photo_url": user.PhotoURL,
			"is_active": user.IsActive,
		})
	}

	meta := domain.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalCount:  total,
		TotalPages:  (total + params.Limit - 1) / params.Limit,
	}

	return SuccessList(c, resp, meta)
}

func (h *SchoolHandler) parseListParams(c *fiber.Ctx) domain.ListParams {
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
			"status":   c.Query("status"),
			"gtk_type": c.Query("gtk_type"),
		},
	}
}

func (h *SchoolHandler) toSchoolResponse(c *fiber.Ctx, school *domain.School) domain.SchoolResponse {
	resp := domain.SchoolResponse{
		ID:        school.ID,
		Name:      school.Name,
		NPSN:      school.NPSN,
		Status:    school.Status,
		Address:   school.Address,
		CreatedAt: school.CreatedAt,
		UpdatedAt: school.UpdatedAt,
	}

	if school.HeadMasterID != nil {
		headMaster, _ := h.schoolService.GetHeadMaster(c.Context(), *school.HeadMasterID)
		if headMaster != nil {
			resp.HeadMaster = &domain.UserRef{
				ID:       headMaster.ID,
				FullName: headMaster.FullName,
				NIP:      headMaster.NIP,
			}
		}
	}

	gtkCount, _ := h.schoolService.CountUsers(c.Context(), school.ID)
	resp.GTKCount = gtkCount

	return resp
}

func (h *SchoolHandler) toSchoolDetailResponse(c *fiber.Ctx, school *domain.School) domain.SchoolDetailResponse {
	resp := domain.SchoolDetailResponse{
		ID:        school.ID,
		Name:      school.Name,
		NPSN:      school.NPSN,
		Status:    school.Status,
		Address:   school.Address,
		CreatedAt: school.CreatedAt,
		UpdatedAt: school.UpdatedAt,
	}

	if school.HeadMasterID != nil {
		headMaster, _ := h.schoolService.GetHeadMaster(c.Context(), *school.HeadMasterID)
		if headMaster != nil {
			resp.HeadMaster = &domain.UserRef{
				ID:       headMaster.ID,
				FullName: headMaster.FullName,
				NIP:      headMaster.NIP,
			}
		}
	}

	resp.GTKCount, _ = h.schoolService.CountUsers(c.Context(), school.ID)
	resp.GuruCount, _ = h.schoolService.CountUsersByType(c.Context(), school.ID, domain.GTKTypeGuru)
	resp.TendikCount, _ = h.schoolService.CountUsersByType(c.Context(), school.ID, domain.GTKTypeTendik)
	resp.KepalaSekolahCount, _ = h.schoolService.CountUsersByType(c.Context(), school.ID, domain.GTKTypeKepalaSekolah)

	return resp
}
