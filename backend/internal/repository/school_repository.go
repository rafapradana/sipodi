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

type SchoolRepository struct {
	db *pgxpool.Pool
}

func NewSchoolRepository(db *pgxpool.Pool) *SchoolRepository {
	return &SchoolRepository{db: db}
}

func (r *SchoolRepository) Create(ctx context.Context, school *domain.School) error {
	query := `
		INSERT INTO schools (id, name, npsn, status, address, head_master_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		school.ID, school.Name, school.NPSN, school.Status, school.Address, school.HeadMasterID,
	).Scan(&school.CreatedAt, &school.UpdatedAt)
}

func (r *SchoolRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.School, error) {
	query := `
		SELECT id, name, npsn, status, address, head_master_id, created_at, updated_at
		FROM schools WHERE id = $1`

	school := &domain.School{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&school.ID, &school.Name, &school.NPSN, &school.Status,
		&school.Address, &school.HeadMasterID, &school.CreatedAt, &school.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return school, err
}

func (r *SchoolRepository) GetByNPSN(ctx context.Context, npsn string) (*domain.School, error) {
	query := `
		SELECT id, name, npsn, status, address, head_master_id, created_at, updated_at
		FROM schools WHERE npsn = $1`

	school := &domain.School{}
	err := r.db.QueryRow(ctx, query, npsn).Scan(
		&school.ID, &school.Name, &school.NPSN, &school.Status,
		&school.Address, &school.HeadMasterID, &school.CreatedAt, &school.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return school, err
}

func (r *SchoolRepository) Update(ctx context.Context, school *domain.School) error {
	query := `
		UPDATE schools SET name = $2, npsn = $3, status = $4, address = $5, head_master_id = $6
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		school.ID, school.Name, school.NPSN, school.Status, school.Address, school.HeadMasterID,
	).Scan(&school.UpdatedAt)
}

func (r *SchoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schools WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *SchoolRepository) List(ctx context.Context, params domain.ListParams) ([]domain.School, int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR npsn ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	if status, ok := params.Filters["status"]; ok && status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM schools %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	orderBy := "created_at DESC"
	if params.Sort != "" {
		if strings.HasPrefix(params.Sort, "-") {
			orderBy = strings.TrimPrefix(params.Sort, "-") + " DESC"
		} else {
			orderBy = params.Sort + " ASC"
		}
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	query := fmt.Sprintf(`
		SELECT id, name, npsn, status, address, head_master_id, created_at, updated_at
		FROM schools %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIndex, argIndex+1,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var schools []domain.School
	for rows.Next() {
		var school domain.School
		err := rows.Scan(
			&school.ID, &school.Name, &school.NPSN, &school.Status,
			&school.Address, &school.HeadMasterID, &school.CreatedAt, &school.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		schools = append(schools, school)
	}

	return schools, total, nil
}

func (r *SchoolRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM schools`
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *SchoolRepository) ExistsByNPSN(ctx context.Context, npsn string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM schools WHERE npsn = $1)`
	err := r.db.QueryRow(ctx, query, npsn).Scan(&exists)
	return exists, err
}

func (r *SchoolRepository) HasUsers(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE school_id = $1)`
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}

func (r *SchoolRepository) GetStatistics(ctx context.Context, params domain.ListParams) ([]domain.SchoolStatistics, int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if status, ok := params.Filters["status"]; ok && status != "" {
		conditions = append(conditions, fmt.Sprintf("s.status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM schools s %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	orderBy := "gtk_count DESC"
	if params.Sort == "talent_count" {
		orderBy = "talent_count DESC"
	} else if params.Sort == "name" {
		orderBy = "s.name ASC"
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	query := fmt.Sprintf(`
		SELECT 
			s.id, s.name, s.npsn, s.status,
			COUNT(DISTINCT u.id) as gtk_count,
			COUNT(DISTINCT t.id) as talent_count,
			COUNT(DISTINCT CASE WHEN t.status = 'pending' THEN t.id END) as pending_count,
			COUNT(DISTINCT CASE WHEN t.status = 'approved' THEN t.id END) as approved_count
		FROM schools s
		LEFT JOIN users u ON u.school_id = s.id AND u.role = 'gtk'
		LEFT JOIN talents t ON t.user_id = u.id
		%s
		GROUP BY s.id, s.name, s.npsn, s.status
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIndex, argIndex+1,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stats []domain.SchoolStatistics
	for rows.Next() {
		var stat domain.SchoolStatistics
		err := rows.Scan(
			&stat.ID, &stat.Name, &stat.NPSN, &stat.Status,
			&stat.GTKCount, &stat.TalentCount, &stat.PendingCount, &stat.ApprovedCount,
		)
		if err != nil {
			return nil, 0, err
		}
		stats = append(stats, stat)
	}

	return stats, total, nil
}
