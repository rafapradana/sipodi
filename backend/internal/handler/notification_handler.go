package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/service"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) List(c *fiber.Ctx) error {
	claims := GetClaims(c)
	params := h.parseListParams(c)

	notifications, total, err := h.notificationService.List(c.Context(), claims.UserID, params)
	if err != nil {
		return InternalError(c)
	}

	var resp []domain.NotificationResponse
	for _, n := range notifications {
		resp = append(resp, domain.NotificationResponse{
			ID:        n.ID,
			Type:      n.Type,
			Message:   n.Message,
			TalentID:  n.TalentID,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt,
		})
	}

	unreadCount, _ := h.notificationService.CountUnread(c.Context(), claims.UserID)

	meta := domain.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalCount:  total,
		TotalPages:  (total + params.Limit - 1) / params.Limit,
	}

	return c.JSON(fiber.Map{
		"data": resp,
		"meta": fiber.Map{
			"current_page": meta.CurrentPage,
			"per_page":     meta.PerPage,
			"total_pages":  meta.TotalPages,
			"total_count":  meta.TotalCount,
			"unread_count": unreadCount,
		},
	})
}

func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	claims := GetClaims(c)
	count, err := h.notificationService.CountUnread(c.Context(), claims.UserID)
	if err != nil {
		return InternalError(c)
	}

	return Success(c, fiber.Map{"unread_count": count})
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "INVALID_ID", "ID tidak valid")
	}

	claims := GetClaims(c)
	err = h.notificationService.MarkAsRead(c.Context(), id, claims.UserID)
	if err != nil {
		switch err {
		case service.ErrNotificationNotFound:
			return NotFound(c, "Notifikasi tidak ditemukan")
		case service.ErrForbidden:
			return Forbidden(c, "Anda tidak memiliki akses ke notifikasi ini")
		default:
			return InternalError(c)
		}
	}

	return Message(c, "Notifikasi ditandai sudah dibaca")
}

func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	claims := GetClaims(c)
	count, err := h.notificationService.MarkAllAsRead(c.Context(), claims.UserID)
	if err != nil {
		return InternalError(c)
	}

	return SuccessWithMessage(c, fiber.Map{"marked_count": count}, "Semua notifikasi ditandai sudah dibaca")
}

func (h *NotificationHandler) parseListParams(c *fiber.Ctx) domain.ListParams {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	return domain.ListParams{
		Page:  page,
		Limit: limit,
		Filters: map[string]string{
			"is_read": c.Query("is_read"),
		},
	}
}
