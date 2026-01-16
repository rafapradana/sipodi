package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/repository"
)

var (
	ErrSchoolNotFound    = errors.New("school not found")
	ErrDuplicateNPSN     = errors.New("npsn already exists")
	ErrSchoolHasUsers    = errors.New("school has users")
	ErrInvalidHeadMaster = errors.New("invalid head master")
)

type SchoolService struct {
	schoolRepo *repository.SchoolRepository
	userRepo   *repository.UserRepository
}

func NewSchoolService(schoolRepo *repository.SchoolRepository, userRepo *repository.UserRepository) *SchoolService {
	return &SchoolService{
		schoolRepo: schoolRepo,
		userRepo:   userRepo,
	}
}

func (s *SchoolService) Create(ctx context.Context, req domain.CreateSchoolRequest) (*domain.School, error) {
	exists, err := s.schoolRepo.ExistsByNPSN(ctx, req.NPSN)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateNPSN
	}

	school := &domain.School{
		ID:      uuid.New(),
		Name:    req.Name,
		NPSN:    req.NPSN,
		Status:  req.Status,
		Address: req.Address,
	}

	if err := s.schoolRepo.Create(ctx, school); err != nil {
		return nil, err
	}

	return school, nil
}

func (s *SchoolService) GetByID(ctx context.Context, id uuid.UUID) (*domain.School, error) {
	school, err := s.schoolRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if school == nil {
		return nil, ErrSchoolNotFound
	}
	return school, nil
}

func (s *SchoolService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateSchoolRequest) (*domain.School, error) {
	school, err := s.schoolRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if school == nil {
		return nil, ErrSchoolNotFound
	}

	if req.Name != nil {
		school.Name = *req.Name
	}
	if req.NPSN != nil {
		// Check if NPSN is unique
		existing, err := s.schoolRepo.GetByNPSN(ctx, *req.NPSN)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrDuplicateNPSN
		}
		school.NPSN = *req.NPSN
	}
	if req.Status != nil {
		school.Status = *req.Status
	}
	if req.Address != nil {
		school.Address = *req.Address
	}
	if req.HeadMasterID != nil {
		// Validate head master
		user, err := s.userRepo.GetByID(ctx, *req.HeadMasterID)
		if err != nil {
			return nil, err
		}
		if user == nil || user.GTKType == nil || *user.GTKType != domain.GTKTypeKepalaSekolah {
			return nil, ErrInvalidHeadMaster
		}
		school.HeadMasterID = req.HeadMasterID
	}

	if err := s.schoolRepo.Update(ctx, school); err != nil {
		return nil, err
	}

	return school, nil
}

func (s *SchoolService) Delete(ctx context.Context, id uuid.UUID) error {
	school, err := s.schoolRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if school == nil {
		return ErrSchoolNotFound
	}

	hasUsers, err := s.schoolRepo.HasUsers(ctx, id)
	if err != nil {
		return err
	}
	if hasUsers {
		return ErrSchoolHasUsers
	}

	return s.schoolRepo.Delete(ctx, id)
}

func (s *SchoolService) List(ctx context.Context, params domain.ListParams) ([]domain.School, int, error) {
	return s.schoolRepo.List(ctx, params)
}

func (s *SchoolService) GetUsers(ctx context.Context, schoolID uuid.UUID, params domain.ListParams) ([]domain.User, int, error) {
	params.Filters["school_id"] = schoolID.String()
	return s.userRepo.List(ctx, params)
}

func (s *SchoolService) GetHeadMaster(ctx context.Context, headMasterID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, headMasterID)
}

func (s *SchoolService) CountUsers(ctx context.Context, schoolID uuid.UUID) (int, error) {
	return s.userRepo.CountBySchool(ctx, schoolID)
}

func (s *SchoolService) CountUsersByType(ctx context.Context, schoolID uuid.UUID, gtkType domain.GTKType) (int, error) {
	return s.userRepo.CountBySchoolAndType(ctx, schoolID, gtkType)
}

func (s *SchoolService) Count(ctx context.Context) (int, error) {
	return s.schoolRepo.Count(ctx)
}
