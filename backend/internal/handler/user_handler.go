package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	claims := GetClaims(c)
	user, err := h.userService.GetByID(c.Context(), claims.UserID)
	if err != nil {
		return InternalError(c)
	}

	resp := h.toUserResponse(c, user)
	return Success(c, resp)
}

func (h *UserHandler) UpdateMe(c *fiber.Ctx) error {
	claims := GetClaims(c)
	var req domain.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	user, err := h.userService.UpdateProfile(c.Context(), claims.UserID, req)
	if err != nil {
		return InternalError(c)
	}

	resp := h.toUserResponse(c, user)
	return SuccessWithMessage(c, resp, "Profil berhasil diperbarui")
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	claims := GetClaims(c)
	var req domain.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	// Validation
	var errors []domain.FieldError
	if req.CurrentPassword == "" {
		errors = append(errors, domain.FieldError{Field: "current_password", Message: "Password lama wajib diisi"})
	}
	if len(req.NewPassword) < 8 {
		errors = append(errors, domain.FieldError{Field: "new_password", Message: "Password minimal 8 karakter"})
	}
	if req.NewPassword != req.NewPasswordConfirmation {
		errors = append(errors, domain.FieldError{Field: "new_password_confirmation", Message: "Konfirmasi password tidak cocok"})
	}
	if len(errors) > 0 {
		return ValidationError(c, errors)
	}

	err := h.userService.ChangePassword(c.Context(), claims.UserID, req)
	if err != nil {
		if err == service.ErrInvalidPassword {
			return BadRequest(c, "INVALID_PASSWORD", "Password lama tidak sesuai")
		}
		return InternalError(c)
	}

	return Message(c, "Password berhasil diubah")
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	params := h.parseListParams(c)
	users, total, err := h.userService.List(c.Context(), params)
	if err != nil {
		return InternalError(c)
	}

	var resp []domain.UserListResponse
	for _, user := range users {
		resp = append(resp, h.toUserListResponse(c, &user))
	}

	meta := domain.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalCount:  total,
		TotalPages:  (total + params.Limit - 1) / params.Limit,
	}

	return SuccessList(c, resp, meta)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	user, err := h.userService.GetByID(c.Context(), id)
	if err != nil {
		if err == service.ErrUserNotFound {
			return NotFound(c, "User tidak ditemukan")
		}
		return InternalError(c)
	}

	resp := h.toUserResponse(c, user)
	return Success(c, resp)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	// Validation
	var errors []domain.FieldError
	if req.Email == "" {
		errors = append(errors, domain.FieldError{Field: "email", Message: "Email wajib diisi"})
	}
	if len(req.Password) < 8 {
		errors = append(errors, domain.FieldError{Field: "password", Message: "Password minimal 8 karakter"})
	}
	if req.FullName == "" {
		errors = append(errors, domain.FieldError{Field: "full_name", Message: "Nama lengkap wajib diisi"})
	}
	if len(errors) > 0 {
		return ValidationError(c, errors)
	}

	// Check authorization for admin_sekolah
	claims := GetClaims(c)
	if claims.Role == domain.RoleAdminSekolah {
		if req.SchoolID == nil || *req.SchoolID != *claims.SchoolID {
			return Forbidden(c, "Anda hanya dapat menambah user di sekolah Anda")
		}
	}

	user, err := h.userService.Create(c.Context(), req)
	if err != nil {
		switch err {
		case service.ErrEmailTaken:
			return Conflict(c, "EMAIL_TAKEN", "Email sudah terdaftar")
		case service.ErrNUPTKTaken:
			return Conflict(c, "NUPTK_TAKEN", "NUPTK sudah terdaftar")
		case service.ErrNIPTaken:
			return Conflict(c, "NIP_TAKEN", "NIP sudah terdaftar")
		default:
			return InternalError(c)
		}
	}

	resp := h.toUserResponse(c, user)
	return SuccessCreated(c, resp, "User berhasil ditambahkan")
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	var req domain.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequest(c, "INVALID_REQUEST", "Request body tidak valid")
	}

	user, err := h.userService.Update(c.Context(), id, req)
	if err != nil {
		if err == service.ErrUserNotFound {
			return NotFound(c, "User tidak ditemukan")
		}
		return InternalError(c)
	}

	resp := h.toUserResponse(c, user)
	return SuccessWithMessage(c, resp, "User berhasil diperbarui")
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	claims := GetClaims(c)
	err = h.userService.Delete(c.Context(), id, claims.UserID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return NotFound(c, "User tidak ditemukan")
		case service.ErrCannotDeleteSelf:
			return BadRequest(c, "CANNOT_DELETE_SELF", "Tidak dapat menghapus akun sendiri")
		default:
			return InternalError(c)
		}
	}

	return SuccessNoContent(c)
}

func (h *UserHandler) Activate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	if err := h.userService.Activate(c.Context(), id); err != nil {
		if err == service.ErrUserNotFound {
			return NotFound(c, "User tidak ditemukan")
		}
		return InternalError(c)
	}

	return Message(c, "User berhasil diaktifkan")
}

func (h *UserHandler) Deactivate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	if err := h.userService.Deactivate(c.Context(), id); err != nil {
		if err == service.ErrUserNotFound {
			return NotFound(c, "User tidak ditemukan")
		}
		return InternalError(c)
	}

	return Message(c, "User berhasil dinonaktifkan")
}

func (h *UserHandler) parseListParams(c *fiber.Ctx) domain.ListParams {
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
			"role":      c.Query("role"),
			"school_id": c.Query("school_id"),
			"gtk_type":  c.Query("gtk_type"),
			"is_active": c.Query("is_active"),
		},
	}
}

func (h *UserHandler) toUserResponse(c *fiber.Ctx, user *domain.User) domain.UserResponse {
	resp := domain.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		FullName:  user.FullName,
		PhotoURL:  user.PhotoURL,
		NUPTK:     user.NUPTK,
		NIP:       user.NIP,
		Gender:    user.Gender,
		GTKType:   user.GTKType,
		Position:  user.Position,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.BirthDate != nil {
		str := user.BirthDate.Format("2006-01-02")
		resp.BirthDate = &str
	}

	if user.SchoolID != nil {
		school, _ := h.userService.GetSchool(c.Context(), *user.SchoolID)
		if school != nil {
			resp.School = &domain.SchoolRef{
				ID:   school.ID,
				Name: school.Name,
				NPSN: school.NPSN,
			}
		}
	}

	return resp
}

func (h *UserHandler) toUserListResponse(c *fiber.Ctx, user *domain.User) domain.UserListResponse {
	resp := domain.UserListResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		FullName:  user.FullName,
		PhotoURL:  user.PhotoURL,
		NUPTK:     user.NUPTK,
		NIP:       user.NIP,
		GTKType:   user.GTKType,
		Position:  user.Position,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	if user.SchoolID != nil {
		school, _ := h.userService.GetSchool(c.Context(), *user.SchoolID)
		if school != nil {
			resp.School = &domain.SchoolRef{
				ID:   school.ID,
				Name: school.Name,
			}
		}
	}

	return resp
}
