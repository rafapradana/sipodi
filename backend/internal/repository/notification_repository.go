package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sipodi/backend/internal/domain"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, talent_id, type, message, is_read)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`

	return r.db.QueryRow(ctx, query,
		notification.ID, notification.UserID, notification.TalentID,
		notification.Type, notification.Message, notification.IsRead,
	).Scan(&notification.CreatedAt)
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	query := `
		SELECT id, user_id, talent_id, type, message, is_read, created_at
		FROM notifications WHERE id = $1`

	notification := &domain.Notification{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&notification.ID, &notification.UserID, &notification.TalentID,
		&notification.Type, &notification.Message, &notification.IsRead,
		&notification.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return notification, err
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false`
	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (r *NotificationRepository) List(ctx context.Context, userID uuid.UUID, params domain.ListParams) ([]domain.Notification, int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if isRead, ok := params.Filters["is_read"]; ok && isRead != "" {
		conditions = append(conditions, fmt.Sprintf("is_read = $%d", argIndex))
		args = append(args, isRead == "true")
		argIndex++
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM notifications %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	query := fmt.Sprintf(`
		SELECT id, user_id, talent_id, type, message, is_read, created_at
		FROM notifications %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var notification domain.Notification
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.TalentID,
			&notification.Type, &notification.Message, &notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, total, nil
}

func (r *NotificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}
