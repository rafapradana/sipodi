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

type TalentRepository struct {
	db *pgxpool.Pool
}

func NewTalentRepository(db *pgxpool.Pool) *TalentRepository {
	return &TalentRepository{db: db}
}

func (r *TalentRepository) Create(ctx context.Context, talent *domain.Talent) error {
	query := `
		INSERT INTO talents (id, user_id, talent_type, status)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		talent.ID, talent.UserID, talent.TalentType, talent.Status,
	).Scan(&talent.CreatedAt, &talent.UpdatedAt)
}

func (r *TalentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Talent, error) {
	query := `
		SELECT id, user_id, talent_type, status, verified_by, verified_at, rejection_reason, created_at, updated_at
		FROM talents WHERE id = $1`

	talent := &domain.Talent{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&talent.ID, &talent.UserID, &talent.TalentType, &talent.Status,
		&talent.VerifiedBy, &talent.VerifiedAt, &talent.RejectionReason,
		&talent.CreatedAt, &talent.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return talent, err
}

func (r *TalentRepository) Update(ctx context.Context, talent *domain.Talent) error {
	query := `
		UPDATE talents SET status = $2, verified_by = $3, verified_at = $4, rejection_reason = $5
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		talent.ID, talent.Status, talent.VerifiedBy, talent.VerifiedAt, talent.RejectionReason,
	).Scan(&talent.UpdatedAt)
}

func (r *TalentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM talents WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *TalentRepository) List(ctx context.Context, params domain.ListParams) ([]domain.Talent, int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if userID, ok := params.Filters["user_id"]; ok && userID != "" {
		conditions = append(conditions, fmt.Sprintf("t.user_id = $%d", argIndex))
		args = append(args, userID)
		argIndex++
	}

	if schoolID, ok := params.Filters["school_id"]; ok && schoolID != "" {
		conditions = append(conditions, fmt.Sprintf("u.school_id = $%d", argIndex))
		args = append(args, schoolID)
		argIndex++
	}

	if talentType, ok := params.Filters["talent_type"]; ok && talentType != "" {
		conditions = append(conditions, fmt.Sprintf("t.talent_type = $%d", argIndex))
		args = append(args, talentType)
		argIndex++
	}

	if status, ok := params.Filters["status"]; ok && status != "" {
		conditions = append(conditions, fmt.Sprintf("t.status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM talents t
		JOIN users u ON t.user_id = u.id
		%s`, whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	orderBy := "t.created_at DESC"
	if params.Sort != "" {
		if strings.HasPrefix(params.Sort, "-") {
			orderBy = "t." + strings.TrimPrefix(params.Sort, "-") + " DESC"
		} else {
			orderBy = "t." + params.Sort + " ASC"
		}
	}

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	query := fmt.Sprintf(`
		SELECT t.id, t.user_id, t.talent_type, t.status, t.verified_by, t.verified_at, t.rejection_reason, t.created_at, t.updated_at
		FROM talents t
		JOIN users u ON t.user_id = u.id
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIndex, argIndex+1,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var talents []domain.Talent
	for rows.Next() {
		var talent domain.Talent
		err := rows.Scan(
			&talent.ID, &talent.UserID, &talent.TalentType, &talent.Status,
			&talent.VerifiedBy, &talent.VerifiedAt, &talent.RejectionReason,
			&talent.CreatedAt, &talent.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		talents = append(talents, talent)
	}

	return talents, total, nil
}

// Training methods
func (r *TalentRepository) CreateTraining(ctx context.Context, training *domain.TalentTraining) error {
	query := `
		INSERT INTO talent_trainings (id, talent_id, activity_name, organizer, start_date, duration_days)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(ctx, query,
		training.ID, training.TalentID, training.ActivityName,
		training.Organizer, training.StartDate, training.DurationDays,
	)
	return err
}

func (r *TalentRepository) GetTrainingByTalentID(ctx context.Context, talentID uuid.UUID) (*domain.TalentTraining, error) {
	query := `
		SELECT id, talent_id, activity_name, organizer, start_date, duration_days
		FROM talent_trainings WHERE talent_id = $1`

	training := &domain.TalentTraining{}
	err := r.db.QueryRow(ctx, query, talentID).Scan(
		&training.ID, &training.TalentID, &training.ActivityName,
		&training.Organizer, &training.StartDate, &training.DurationDays,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return training, err
}

func (r *TalentRepository) UpdateTraining(ctx context.Context, training *domain.TalentTraining) error {
	query := `
		UPDATE talent_trainings SET activity_name = $2, organizer = $3, start_date = $4, duration_days = $5
		WHERE talent_id = $1`

	_, err := r.db.Exec(ctx, query,
		training.TalentID, training.ActivityName, training.Organizer,
		training.StartDate, training.DurationDays,
	)
	return err
}

// Mentor methods
func (r *TalentRepository) CreateMentor(ctx context.Context, mentor *domain.TalentCompetitionMentor) error {
	query := `
		INSERT INTO talent_competition_mentors (id, talent_id, competition_name, level, organizer, field, achievement, certificate_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(ctx, query,
		mentor.ID, mentor.TalentID, mentor.CompetitionName, mentor.Level,
		mentor.Organizer, mentor.Field, mentor.Achievement, mentor.CertificateURL,
	)
	return err
}

func (r *TalentRepository) GetMentorByTalentID(ctx context.Context, talentID uuid.UUID) (*domain.TalentCompetitionMentor, error) {
	query := `
		SELECT id, talent_id, competition_name, level, organizer, field, achievement, certificate_url
		FROM talent_competition_mentors WHERE talent_id = $1`

	mentor := &domain.TalentCompetitionMentor{}
	err := r.db.QueryRow(ctx, query, talentID).Scan(
		&mentor.ID, &mentor.TalentID, &mentor.CompetitionName, &mentor.Level,
		&mentor.Organizer, &mentor.Field, &mentor.Achievement, &mentor.CertificateURL,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return mentor, err
}

func (r *TalentRepository) UpdateMentor(ctx context.Context, mentor *domain.TalentCompetitionMentor) error {
	query := `
		UPDATE talent_competition_mentors SET competition_name = $2, level = $3, organizer = $4, field = $5, achievement = $6, certificate_url = $7
		WHERE talent_id = $1`

	_, err := r.db.Exec(ctx, query,
		mentor.TalentID, mentor.CompetitionName, mentor.Level,
		mentor.Organizer, mentor.Field, mentor.Achievement, mentor.CertificateURL,
	)
	return err
}

// Participant methods
func (r *TalentRepository) CreateParticipant(ctx context.Context, participant *domain.TalentCompetitionParticipant) error {
	query := `
		INSERT INTO talent_competition_participants (id, talent_id, competition_name, level, organizer, field, start_date, duration_days, competition_field, achievement, certificate_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.Exec(ctx, query,
		participant.ID, participant.TalentID, participant.CompetitionName, participant.Level,
		participant.Organizer, participant.Field, participant.StartDate, participant.DurationDays,
		participant.CompetitionField, participant.Achievement, participant.CertificateURL,
	)
	return err
}

func (r *TalentRepository) GetParticipantByTalentID(ctx context.Context, talentID uuid.UUID) (*domain.TalentCompetitionParticipant, error) {
	query := `
		SELECT id, talent_id, competition_name, level, organizer, field, start_date, duration_days, competition_field, achievement, certificate_url
		FROM talent_competition_participants WHERE talent_id = $1`

	participant := &domain.TalentCompetitionParticipant{}
	err := r.db.QueryRow(ctx, query, talentID).Scan(
		&participant.ID, &participant.TalentID, &participant.CompetitionName, &participant.Level,
		&participant.Organizer, &participant.Field, &participant.StartDate, &participant.DurationDays,
		&participant.CompetitionField, &participant.Achievement, &participant.CertificateURL,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return participant, err
}

func (r *TalentRepository) UpdateParticipant(ctx context.Context, participant *domain.TalentCompetitionParticipant) error {
	query := `
		UPDATE talent_competition_participants SET competition_name = $2, level = $3, organizer = $4, field = $5, start_date = $6, duration_days = $7, competition_field = $8, achievement = $9, certificate_url = $10
		WHERE talent_id = $1`

	_, err := r.db.Exec(ctx, query,
		participant.TalentID, participant.CompetitionName, participant.Level,
		participant.Organizer, participant.Field, participant.StartDate, participant.DurationDays,
		participant.CompetitionField, participant.Achievement, participant.CertificateURL,
	)
	return err
}

// Interest methods
func (r *TalentRepository) CreateInterest(ctx context.Context, interest *domain.TalentInterest) error {
	query := `
		INSERT INTO talent_interests (id, talent_id, interest_name, description, certificate_url)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query,
		interest.ID, interest.TalentID, interest.InterestName,
		interest.Description, interest.CertificateURL,
	)
	return err
}

func (r *TalentRepository) GetInterestByTalentID(ctx context.Context, talentID uuid.UUID) (*domain.TalentInterest, error) {
	query := `
		SELECT id, talent_id, interest_name, description, certificate_url
		FROM talent_interests WHERE talent_id = $1`

	interest := &domain.TalentInterest{}
	err := r.db.QueryRow(ctx, query, talentID).Scan(
		&interest.ID, &interest.TalentID, &interest.InterestName,
		&interest.Description, &interest.CertificateURL,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return interest, err
}

func (r *TalentRepository) UpdateInterest(ctx context.Context, interest *domain.TalentInterest) error {
	query := `
		UPDATE talent_interests SET interest_name = $2, description = $3, certificate_url = $4
		WHERE talent_id = $1`

	_, err := r.db.Exec(ctx, query,
		interest.TalentID, interest.InterestName, interest.Description, interest.CertificateURL,
	)
	return err
}

// Statistics methods
func (r *TalentRepository) CountByStatus(ctx context.Context) (map[string]int, error) {
	query := `SELECT status, COUNT(*) FROM talents GROUP BY status`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
	}
	return result, nil
}

func (r *TalentRepository) CountByType(ctx context.Context) (map[string]int, error) {
	query := `SELECT talent_type, COUNT(*) FROM talents GROUP BY talent_type`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var talentType string
		var count int
		if err := rows.Scan(&talentType, &count); err != nil {
			return nil, err
		}
		result[talentType] = count
	}
	return result, nil
}

func (r *TalentRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (map[string]int, error) {
	query := `SELECT status, COUNT(*) FROM talents WHERE user_id = $1 GROUP BY status`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{"total": 0, "pending": 0, "approved": 0, "rejected": 0}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
		result["total"] += count
	}
	return result, nil
}

func (r *TalentRepository) CountBySchoolID(ctx context.Context, schoolID uuid.UUID) (map[string]int, error) {
	query := `
		SELECT t.status, COUNT(*) FROM talents t
		JOIN users u ON t.user_id = u.id
		WHERE u.school_id = $1
		GROUP BY t.status`
	rows, err := r.db.Query(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{"total": 0, "pending": 0, "approved": 0, "rejected": 0}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
		result["total"] += count
	}
	return result, nil
}

func (r *TalentRepository) CountPendingBySchoolID(ctx context.Context, schoolID uuid.UUID) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM talents t
		JOIN users u ON t.user_id = u.id
		WHERE u.school_id = $1 AND t.status = 'pending'`
	err := r.db.QueryRow(ctx, query, schoolID).Scan(&count)
	return count, err
}

func (r *TalentRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM talents`
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *TalentRepository) GetRecentByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Talent, error) {
	query := `
		SELECT id, user_id, talent_type, status, verified_by, verified_at, rejection_reason, created_at, updated_at
		FROM talents WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var talents []domain.Talent
	for rows.Next() {
		var talent domain.Talent
		err := rows.Scan(
			&talent.ID, &talent.UserID, &talent.TalentType, &talent.Status,
			&talent.VerifiedBy, &talent.VerifiedAt, &talent.RejectionReason,
			&talent.CreatedAt, &talent.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		talents = append(talents, talent)
	}
	return talents, nil
}

func (r *TalentRepository) ResetStatus(ctx context.Context, talentID uuid.UUID) error {
	query := `
		UPDATE talents SET status = 'pending', verified_by = NULL, verified_at = NULL, rejection_reason = NULL
		WHERE id = $1`
	_, err := r.db.Exec(ctx, query, talentID)
	return err
}

func (r *TalentRepository) GetStatistics(ctx context.Context, groupBy string, schoolID *string, dateFrom, dateTo string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var conditions []string
	var args []interface{}
	argIndex := 1

	if schoolID != nil && *schoolID != "" {
		conditions = append(conditions, fmt.Sprintf("u.school_id = $%d", argIndex))
		args = append(args, *schoolID)
		argIndex++
	}

	if dateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= $%d", argIndex))
		args = append(args, dateFrom)
		argIndex++
	}

	if dateTo != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= $%d", argIndex))
		args = append(args, dateTo)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	switch groupBy {
	case "type":
		query := fmt.Sprintf(`
			SELECT t.talent_type, COUNT(*) 
			FROM talents t
			JOIN users u ON t.user_id = u.id
			%s
			GROUP BY t.talent_type`, whereClause)
		rows, err := r.db.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		data := make(map[string]int)
		for rows.Next() {
			var talentType string
			var count int
			if err := rows.Scan(&talentType, &count); err != nil {
				return nil, err
			}
			data[talentType] = count
		}
		result["by_type"] = data

	case "status":
		query := fmt.Sprintf(`
			SELECT t.status, COUNT(*) 
			FROM talents t
			JOIN users u ON t.user_id = u.id
			%s
			GROUP BY t.status`, whereClause)
		rows, err := r.db.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		data := make(map[string]int)
		for rows.Next() {
			var status string
			var count int
			if err := rows.Scan(&status, &count); err != nil {
				return nil, err
			}
			data[status] = count
		}
		result["by_status"] = data

	case "level":
		// Get from competition mentors and participants
		query := fmt.Sprintf(`
			SELECT level, COUNT(*) FROM (
				SELECT tcm.level FROM talent_competition_mentors tcm
				JOIN talents t ON tcm.talent_id = t.id
				JOIN users u ON t.user_id = u.id
				%s
				UNION ALL
				SELECT tcp.level FROM talent_competition_participants tcp
				JOIN talents t ON tcp.talent_id = t.id
				JOIN users u ON t.user_id = u.id
				%s
			) combined
			GROUP BY level`, whereClause, whereClause)

		// Double the args for UNION ALL
		allArgs := append(args, args...)
		rows, err := r.db.Query(ctx, query, allArgs...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		data := make(map[string]int)
		for rows.Next() {
			var level string
			var count int
			if err := rows.Scan(&level, &count); err != nil {
				return nil, err
			}
			data[level] = count
		}
		result["by_level"] = data

	case "field":
		query := fmt.Sprintf(`
			SELECT field, COUNT(*) FROM (
				SELECT tcm.field FROM talent_competition_mentors tcm
				JOIN talents t ON tcm.talent_id = t.id
				JOIN users u ON t.user_id = u.id
				%s
				UNION ALL
				SELECT tcp.field FROM talent_competition_participants tcp
				JOIN talents t ON tcp.talent_id = t.id
				JOIN users u ON t.user_id = u.id
				%s
			) combined
			GROUP BY field`, whereClause, whereClause)

		allArgs := append(args, args...)
		rows, err := r.db.Query(ctx, query, allArgs...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		data := make(map[string]int)
		for rows.Next() {
			var field string
			var count int
			if err := rows.Scan(&field, &count); err != nil {
				return nil, err
			}
			data[field] = count
		}
		result["by_field"] = data
	}

	return result, nil
}
