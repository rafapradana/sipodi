package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/service"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	claims := GetClaims(c)

	var summary *domain.DashboardSummary
	var err error

	switch claims.Role {
	case domain.RoleSuperAdmin:
		summary, err = h.dashboardService.GetSuperAdminSummary(c.Context())
	case domain.RoleAdminSekolah:
		if claims.SchoolID != nil {
			summary, err = h.dashboardService.GetAdminSekolahSummary(c.Context(), *claims.SchoolID)
		}
	case domain.RoleGTK:
		summary, err = h.dashboardService.GetGTKSummary(c.Context(), claims.UserID)
	}

	if err != nil {
		return InternalError(c)
	}

	return Success(c, summary)
}

func (h *DashboardHandler) GetSchoolsStatistics(c *fiber.Ctx) error {
	params := domain.ListParams{
		Page:  1,
		Limit: 20,
		Sort:  c.Query("sort", "gtk_count"),
		Filters: map[string]string{
			"status": c.Query("status"),
		},
	}

	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil {
		params.Page = page
	}
	if limit, err := strconv.Atoi(c.Query("limit", "20")); err == nil {
		params.Limit = limit
	}

	stats, total, err := h.dashboardService.GetSchoolsStatistics(c.Context(), params)
	if err != nil {
		return InternalError(c)
	}

	meta := domain.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalCount:  total,
		TotalPages:  (total + params.Limit - 1) / params.Limit,
	}

	return SuccessList(c, stats, meta)
}

func (h *DashboardHandler) GetTalentsStatistics(c *fiber.Ctx) error {
	claims := GetClaims(c)
	groupBy := c.Query("group_by", "type")

	var schoolID *string
	if claims.Role == domain.RoleAdminSekolah && claims.SchoolID != nil {
		sid := claims.SchoolID.String()
		schoolID = &sid
	} else if c.Query("school_id") != "" {
		sid := c.Query("school_id")
		schoolID = &sid
	}

	stats, err := h.dashboardService.GetTalentsStatistics(c.Context(), groupBy, schoolID, c.Query("date_from"), c.Query("date_to"))
	if err != nil {
		return InternalError(c)
	}

	return Success(c, stats)
}
