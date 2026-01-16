package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/repository"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationService(notificationRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *NotificationService) List(ctx context.Context, userID uuid.UUID, params domain.ListParams) ([]domain.Notification, int, error) {
	return s.notificationRepo.List(ctx, userID, params)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if notification == nil {
		return ErrNotificationNotFound
	}
	if notification.UserID != userID {
		return ErrForbidden
	}

	return s.notificationRepo.MarkAsRead(ctx, id)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.notificationRepo.CountUnread(ctx, userID)
}
