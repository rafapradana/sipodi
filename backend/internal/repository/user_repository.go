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

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, full_name, photo_url, nuptk, nip, gender, birth_date, gtk_type, position, school_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Role, user.FullName,
		user.PhotoURL, user.NUPTK, user.NIP, user.Gender, user.BirthDate,
		user.GTKType, user.Position, user.SchoolID, user.IsActive,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, full_name, photo_url, nuptk, nip, gender, birth_date, gtk_type, position, school_id, is_active, created_at, updated_at
		FROM users WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.FullName,
		&user.PhotoURL, &user.NUPTK, &user.NIP, &user.Gender, &user.BirthDate,
		&user.GTKType, &user.Position, &user.SchoolID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, full_name, photo_url, nuptk, nip, gender, birth_date, gtk_type, position, school_id, is_active, created_at, updated_at
		FROM users WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.FullName,
		&user.PhotoURL, &user.NUPTK, &user.NIP, &user.Gender, &user.BirthDate,
		&user.GTKType, &user.Position, &user.SchoolID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET
			full_name = $2, photo_url = $3, nuptk = $4, nip = $5, gender = $6,
			birth_date = $7, gtk_type = $8, position = $9, school_id = $10, is_active = $11
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		user.ID, user.FullName, user.PhotoURL, user.NUPTK, user.NIP,
		user.Gender, user.BirthDate, user.GTKType, user.Position,
		user.SchoolID, user.IsActive,
	).Scan(&user.UpdatedAt)
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, passwordHash)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *UserRepository) List(ctx context.Context, params domain.ListParams) ([]domain.User, int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(full_name ILIKE $%d OR email ILIKE $%d OR nuptk ILIKE $%d OR nip ILIKE $%d)",
			argIndex, argIndex, argIndex, argIndex,
		))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	if role, ok := params.Filters["role"]; ok && role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, role)
		argIndex++
	}

	if schoolID, ok := params.Filters["school_id"]; ok && schoolID != "" {
		conditions = append(conditions, fmt.Sprintf("school_id = $%d", argIndex))
		args = append(args, schoolID)
		argIndex++
	}

	if gtkType, ok := params.Filters["gtk_type"]; ok && gtkType != "" {
		conditions = append(conditions, fmt.Sprintf("gtk_type = $%d", argIndex))
		args = append(args, gtkType)
		argIndex++
	}

	if isActive, ok := params.Filters["is_active"]; ok && isActive != "" {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, isActive == "true")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Sort
	orderBy := "created_at DESC"
	if params.Sort != "" {
		if strings.HasPrefix(params.Sort, "-") {
			orderBy = strings.TrimPrefix(params.Sort, "-") + " DESC"
		} else {
			orderBy = params.Sort + " ASC"
		}
	}

	// Pagination
	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	query := fmt.Sprintf(`
		SELECT id, email, password_hash, role, full_name, photo_url, nuptk, nip, gender, birth_date, gtk_type, position, school_id, is_active, created_at, updated_at
		FROM users %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIndex, argIndex+1,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.FullName,
			&user.PhotoURL, &user.NUPTK, &user.NIP, &user.Gender, &user.BirthDate,
			&user.GTKType, &user.Position, &user.SchoolID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) ExistsByNUPTK(ctx context.Context, nuptk string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE nuptk = $1)`
	err := r.db.QueryRow(ctx, query, nuptk).Scan(&exists)
	return exists, err
}

func (r *UserRepository) ExistsByNIP(ctx context.Context, nip string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE nip = $1)`
	err := r.db.QueryRow(ctx, query, nip).Scan(&exists)
	return exists, err
}

func (r *UserRepository) CountBySchool(ctx context.Context, schoolID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE school_id = $1`
	err := r.db.QueryRow(ctx, query, schoolID).Scan(&count)
	return count, err
}

func (r *UserRepository) CountBySchoolAndType(ctx context.Context, schoolID uuid.UUID, gtkType domain.GTKType) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE school_id = $1 AND gtk_type = $2`
	err := r.db.QueryRow(ctx, query, schoolID, gtkType).Scan(&count)
	return count, err
}

func (r *UserRepository) CountByRole(ctx context.Context, role domain.UserRole) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE role = $1`
	err := r.db.QueryRow(ctx, query, role).Scan(&count)
	return count, err
}

func (r *UserRepository) CountByGTKType(ctx context.Context) (map[string]int, error) {
	query := `SELECT gtk_type, COUNT(*) FROM users WHERE gtk_type IS NOT NULL GROUP BY gtk_type`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var gtkType string
		var count int
		if err := rows.Scan(&gtkType, &count); err != nil {
			return nil, err
		}
		result[gtkType] = count
	}
	return result, nil
}
