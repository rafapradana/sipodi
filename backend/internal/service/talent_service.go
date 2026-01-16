package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/repository"
)

var (
	ErrTalentNotFound  = errors.New("talent not found")
	ErrAlreadyVerified = errors.New("talent already verified")
	ErrForbidden       = errors.New("forbidden")
)

type TalentService struct {
	talentRepo       *repository.TalentRepository
	userRepo         *repository.UserRepository
	notificationRepo *repository.NotificationRepository
}

func NewTalentService(
	talentRepo *repository.TalentRepository,
	userRepo *repository.UserRepository,
	notificationRepo *repository.NotificationRepository,
) *TalentService {
	return &TalentService{
		talentRepo:       talentRepo,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *TalentService) Create(ctx context.Context, userID uuid.UUID, req domain.CreateTalentRequest, certificateURL *string) (*domain.Talent, error) {
	talent := &domain.Talent{
		ID:         uuid.New(),
		UserID:     userID,
		TalentType: req.TalentType,
		Status:     domain.TalentStatusPending,
	}

	if err := s.talentRepo.Create(ctx, talent); err != nil {
		return nil, err
	}

	// Create detail based on type
	switch req.TalentType {
	case domain.TalentTypePesertaPelatihan:
		if err := s.createTrainingDetail(ctx, talent.ID, req.Detail); err != nil {
			return nil, err
		}
	case domain.TalentTypePembimbingLomba:
		if err := s.createMentorDetail(ctx, talent.ID, req.Detail, certificateURL); err != nil {
			return nil, err
		}
	case domain.TalentTypePesertaLomba:
		if err := s.createParticipantDetail(ctx, talent.ID, req.Detail, certificateURL); err != nil {
			return nil, err
		}
	case domain.TalentTypeMinatBakat:
		if err := s.createInterestDetail(ctx, talent.ID, req.Detail, certificateURL); err != nil {
			return nil, err
		}
	}

	return talent, nil
}

func (s *TalentService) createTrainingDetail(ctx context.Context, talentID uuid.UUID, detail interface{}) error {
	detailBytes, _ := json.Marshal(detail)
	var d domain.TrainingDetail
	json.Unmarshal(detailBytes, &d)

	startDate, _ := time.Parse("2006-01-02", d.StartDate)
	training := &domain.TalentTraining{
		ID:           uuid.New(),
		TalentID:     talentID,
		ActivityName: d.ActivityName,
		Organizer:    d.Organizer,
		StartDate:    startDate,
		DurationDays: d.DurationDays,
	}
	return s.talentRepo.CreateTraining(ctx, training)
}

func (s *TalentService) createMentorDetail(ctx context.Context, talentID uuid.UUID, detail interface{}, certURL *string) error {
	detailBytes, _ := json.Marshal(detail)
	var d domain.MentorDetail
	json.Unmarshal(detailBytes, &d)

	mentor := &domain.TalentCompetitionMentor{
		ID:              uuid.New(),
		TalentID:        talentID,
		CompetitionName: d.CompetitionName,
		Level:           d.Level,
		Organizer:       d.Organizer,
		Field:           d.Field,
		Achievement:     d.Achievement,
		CertificateURL:  certURL,
	}
	return s.talentRepo.CreateMentor(ctx, mentor)
}

func (s *TalentService) createParticipantDetail(ctx context.Context, talentID uuid.UUID, detail interface{}, certURL *string) error {
	detailBytes, _ := json.Marshal(detail)
	var d domain.ParticipantDetail
	json.Unmarshal(detailBytes, &d)

	startDate, _ := time.Parse("2006-01-02", d.StartDate)
	participant := &domain.TalentCompetitionParticipant{
		ID:               uuid.New(),
		TalentID:         talentID,
		CompetitionName:  d.CompetitionName,
		Level:            d.Level,
		Organizer:        d.Organizer,
		Field:            d.Field,
		StartDate:        startDate,
		DurationDays:     d.DurationDays,
		CompetitionField: d.CompetitionField,
		Achievement:      d.Achievement,
		CertificateURL:   certURL,
	}
	return s.talentRepo.CreateParticipant(ctx, participant)
}

func (s *TalentService) createInterestDetail(ctx context.Context, talentID uuid.UUID, detail interface{}, certURL *string) error {
	detailBytes, _ := json.Marshal(detail)
	var d domain.InterestDetail
	json.Unmarshal(detailBytes, &d)

	interest := &domain.TalentInterest{
		ID:             uuid.New(),
		TalentID:       talentID,
		InterestName:   d.InterestName,
		Description:    d.Description,
		CertificateURL: certURL,
	}
	return s.talentRepo.CreateInterest(ctx, interest)
}

func (s *TalentService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Talent, error) {
	talent, err := s.talentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if talent == nil {
		return nil, ErrTalentNotFound
	}
	return talent, nil
}

func (s *TalentService) GetDetail(ctx context.Context, talent *domain.Talent) (interface{}, *string, error) {
	switch talent.TalentType {
	case domain.TalentTypePesertaPelatihan:
		detail, err := s.talentRepo.GetTrainingByTalentID(ctx, talent.ID)
		if err != nil {
			return nil, nil, err
		}
		return detail, nil, nil
	case domain.TalentTypePembimbingLomba:
		detail, err := s.talentRepo.GetMentorByTalentID(ctx, talent.ID)
		if err != nil {
			return nil, nil, err
		}
		return detail, detail.CertificateURL, nil
	case domain.TalentTypePesertaLomba:
		detail, err := s.talentRepo.GetParticipantByTalentID(ctx, talent.ID)
		if err != nil {
			return nil, nil, err
		}
		return detail, detail.CertificateURL, nil
	case domain.TalentTypeMinatBakat:
		detail, err := s.talentRepo.GetInterestByTalentID(ctx, talent.ID)
		if err != nil {
			return nil, nil, err
		}
		return detail, detail.CertificateURL, nil
	}
	return nil, nil, nil
}

func (s *TalentService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req domain.UpdateTalentRequest, certificateURL *string) (*domain.Talent, error) {
	talent, err := s.talentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if talent == nil {
		return nil, ErrTalentNotFound
	}
	if talent.UserID != userID {
		return nil, ErrForbidden
	}

	// Update detail based on type
	switch talent.TalentType {
	case domain.TalentTypePesertaPelatihan:
		if err := s.updateTrainingDetail(ctx, talent.ID, req.Detail); err != nil {
			return nil, err
		}
	case domain.TalentTypePembimbingLomba:
		if err := s.updateMentorDetail(ctx, talent.ID, req.Detail, certificateURL); err != nil {
			return nil, err
		}
	case domain.TalentTypePesertaLomba:
		if err := s.updateParticipantDetail(ctx, talent.ID, req.Detail, certificateURL); err != nil {
			return nil, err
		}
	case domain.TalentTypeMinatBakat:
		if err := s.updateInterestDetail(ctx, talent.ID, req.Detail, certificateURL); err != nil {
			return nil, err
		}
	}

	// Reset status to pending
	if err := s.talentRepo.ResetStatus(ctx, talent.ID); err != nil {
		return nil, err
	}

	return s.talentRepo.GetByID(ctx, id)
}

func (s *TalentService) updateTrainingDetail(ctx context.Context, talentID uuid.UUID, detail interface{}) error {
	detailBytes, _ := json.Marshal(detail)
	var d domain.TrainingDetail
	json.Unmarshal(detailBytes, &d)

	startDate, _ := time.Parse("2006-01-02", d.StartDate)
	training := &domain.TalentTraining{
		TalentID:     talentID,
		ActivityName: d.ActivityName,
		Organizer:    d.Organizer,
		StartDate:    startDate,
		DurationDays: d.DurationDays,
	}
	return s.talentRepo.UpdateTraining(ctx, training)
}

func (s *TalentService) updateMentorDetail(ctx context.Context, talentID uuid.UUID, detail interface{}, certURL *string) error {
	detailBytes, _ := json.Marshal(detail)
	var d domain.MentorDetail
	json.Unmarshal(detailBytes, &d)

	mentor := &domain.TalentCompetitionMentor{
		TalentID:        talentID,
		CompetitionName: d.CompetitionName,
		Level:           d.Level,
		Organizer:       d.Organizer,
		Field:           d.Field,
		Achievement:     d.Achievement,
		CertificateURL:  certURL,
	}
	return s.talentRepo.UpdateMentor(ctx, mentor)
}

func (s *TalentService) updateParticipantDetail(ctx context.Context, talentID uuid.UUID, detail interface{}, certURL *string) error {
	detailBytes, _ := json.Marshal(detail)
	var d domain.ParticipantDetail
	json.Unmarshal(detailBytes, &d)

	startDate, _ := time.Parse("2006-01-02", d.StartDate)
	participant := &domain.TalentCompetitionParticipant{
		TalentID:         talentID,
		CompetitionName:  d.CompetitionName,
		Level:            d.Level,
		Organizer:        d.Organizer,
		Field:            d.Field,
		StartDate:        startDate,
		DurationDays:     d.DurationDays,
		CompetitionField: d.CompetitionField,
		Achievement:      d.Achievement,
		CertificateURL:   certURL,
	}
	return s.talentRepo.UpdateParticipant(ctx, participant)
}

func (s *TalentService) updateInterestDetail(ctx context.Context, talentID uuid.UUID, detail interface{}, certURL *string) error {
	detailBytes, _ := json.Marshal(detail)
	var d domain.InterestDetail
	json.Unmarshal(detailBytes, &d)

	interest := &domain.TalentInterest{
		TalentID:       talentID,
		InterestName:   d.InterestName,
		Description:    d.Description,
		CertificateURL: certURL,
	}
	return s.talentRepo.UpdateInterest(ctx, interest)
}

func (s *TalentService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	talent, err := s.talentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if talent == nil {
		return ErrTalentNotFound
	}
	if talent.UserID != userID {
		return ErrForbidden
	}

	return s.talentRepo.Delete(ctx, id)
}

func (s *TalentService) List(ctx context.Context, params domain.ListParams) ([]domain.Talent, int, error) {
	return s.talentRepo.List(ctx, params)
}

func (s *TalentService) GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// Verification methods
func (s *TalentService) Approve(ctx context.Context, id uuid.UUID, verifierID uuid.UUID) (*domain.Talent, error) {
	talent, err := s.talentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if talent == nil {
		return nil, ErrTalentNotFound
	}
	if talent.Status != domain.TalentStatusPending {
		return nil, ErrAlreadyVerified
	}

	now := time.Now()
	talent.Status = domain.TalentStatusApproved
	talent.VerifiedBy = &verifierID
	talent.VerifiedAt = &now

	if err := s.talentRepo.Update(ctx, talent); err != nil {
		return nil, err
	}

	// Create notification
	s.createNotification(ctx, talent.UserID, talent.ID, domain.NotificationTalentApproved, "Talenta Anda telah disetujui")

	return talent, nil
}

func (s *TalentService) Reject(ctx context.Context, id uuid.UUID, verifierID uuid.UUID, reason string) (*domain.Talent, error) {
	talent, err := s.talentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if talent == nil {
		return nil, ErrTalentNotFound
	}
	if talent.Status != domain.TalentStatusPending {
		return nil, ErrAlreadyVerified
	}

	now := time.Now()
	talent.Status = domain.TalentStatusRejected
	talent.VerifiedBy = &verifierID
	talent.VerifiedAt = &now
	talent.RejectionReason = &reason

	if err := s.talentRepo.Update(ctx, talent); err != nil {
		return nil, err
	}

	// Create notification
	message := "Talenta Anda ditolak. Alasan: " + reason
	s.createNotification(ctx, talent.UserID, talent.ID, domain.NotificationTalentRejected, message)

	return talent, nil
}

func (s *TalentService) BatchApprove(ctx context.Context, ids []uuid.UUID, verifierID uuid.UUID) (*domain.BatchResult, error) {
	result := &domain.BatchResult{
		FailedIDs: []domain.FailedItem{},
	}

	for _, id := range ids {
		_, err := s.Approve(ctx, id, verifierID)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, domain.FailedItem{
				ID:     id,
				Reason: err.Error(),
			})
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

func (s *TalentService) BatchReject(ctx context.Context, ids []uuid.UUID, verifierID uuid.UUID, reason string) (*domain.BatchResult, error) {
	result := &domain.BatchResult{
		FailedIDs: []domain.FailedItem{},
	}

	for _, id := range ids {
		_, err := s.Reject(ctx, id, verifierID, reason)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, domain.FailedItem{
				ID:     id,
				Reason: err.Error(),
			})
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

func (s *TalentService) createNotification(ctx context.Context, userID uuid.UUID, talentID uuid.UUID, notifType domain.NotificationType, message string) {
	notification := &domain.Notification{
		ID:       uuid.New(),
		UserID:   userID,
		TalentID: &talentID,
		Type:     notifType,
		Message:  message,
		IsRead:   false,
	}
	s.notificationRepo.Create(ctx, notification)
}

// Statistics
func (s *TalentService) CountByStatus(ctx context.Context) (map[string]int, error) {
	return s.talentRepo.CountByStatus(ctx)
}

func (s *TalentService) CountByType(ctx context.Context) (map[string]int, error) {
	return s.talentRepo.CountByType(ctx)
}

func (s *TalentService) CountByUserID(ctx context.Context, userID uuid.UUID) (map[string]int, error) {
	return s.talentRepo.CountByUserID(ctx, userID)
}

func (s *TalentService) CountBySchoolID(ctx context.Context, schoolID uuid.UUID) (map[string]int, error) {
	return s.talentRepo.CountBySchoolID(ctx, schoolID)
}

func (s *TalentService) CountPendingBySchoolID(ctx context.Context, schoolID uuid.UUID) (int, error) {
	return s.talentRepo.CountPendingBySchoolID(ctx, schoolID)
}

func (s *TalentService) Count(ctx context.Context) (int, error) {
	return s.talentRepo.Count(ctx)
}

func (s *TalentService) GetRecentByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Talent, error) {
	return s.talentRepo.GetRecentByUserID(ctx, userID, limit)
}
