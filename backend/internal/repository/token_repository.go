package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sipodi/backend/internal/domain"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`

	return r.db.QueryRow(ctx, query,
		token.ID, token.UserID, token.TokenHash, token.ExpiresAt,
	).Scan(&token.CreatedAt)
}

func (r *TokenRepository) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND expires_at > $2`

	token := &domain.RefreshToken{}
	err := r.db.QueryRow(ctx, query, tokenHash, time.Now()).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return token, err
}

func (r *TokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *TokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (r *TokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`
	_, err := r.db.Exec(ctx, query, time.Now())
	return err
}

func (r *TokenRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM refresh_tokens WHERE user_id = $1 AND expires_at > $2`
	err := r.db.QueryRow(ctx, query, userID, time.Now()).Scan(&count)
	return count, err
}
