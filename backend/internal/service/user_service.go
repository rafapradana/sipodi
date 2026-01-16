package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/repository"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailTaken       = errors.New("email already taken")
	ErrNUPTKTaken       = errors.New("nuptk already taken")
	ErrNIPTaken         = errors.New("nip already taken")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrCannotDeleteSelf = errors.New("cannot delete self")
)

type UserService struct {
	userRepo   *repository.UserRepository
	schoolRepo *repository.SchoolRepository
}

func NewUserService(userRepo *repository.UserRepository, schoolRepo *repository.SchoolRepository) *UserService {
	return &UserService{
		userRepo:   userRepo,
		schoolRepo: schoolRepo,
	}
}

func (s *UserService) Create(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	// Check email
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailTaken
	}

	// Check NUPTK
	if req.NUPTK != nil && *req.NUPTK != "" {
		exists, err := s.userRepo.ExistsByNUPTK(ctx, *req.NUPTK)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrNUPTKTaken
		}
	}

	// Check NIP
	if req.NIP != nil && *req.NIP != "" {
		exists, err := s.userRepo.ExistsByNIP(ctx, *req.NIP)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrNIPTaken
		}
	}

	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var birthDate *time.Time
	if req.BirthDate != nil && *req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err == nil {
			birthDate = &t
		}
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
		FullName:     req.FullName,
		NUPTK:        req.NUPTK,
		NIP:          req.NIP,
		Gender:       req.Gender,
		BirthDate:    birthDate,
		GTKType:      req.GTKType,
		Position:     req.Position,
		SchoolID:     req.SchoolID,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.NUPTK != nil {
		user.NUPTK = req.NUPTK
	}
	if req.NIP != nil {
		user.NIP = req.NIP
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.BirthDate != nil && *req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err == nil {
			user.BirthDate = &t
		}
	}
	if req.GTKType != nil {
		user.GTKType = req.GTKType
	}
	if req.Position != nil {
		user.Position = req.Position
	}
	if req.SchoolID != nil {
		user.SchoolID = req.SchoolID
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, req domain.UpdateProfileRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.BirthDate != nil && *req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err == nil {
			user.BirthDate = &t
		}
	}
	if req.Position != nil {
		user.Position = req.Position
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, req domain.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if !CheckPassword(req.CurrentPassword, user.PasswordHash) {
		return ErrInvalidPassword
	}

	newHash, err := HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, id, newHash)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) error {
	if id == currentUserID {
		return ErrCannotDeleteSelf
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) Activate(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.IsActive = true
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) Deactivate(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.IsActive = false
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) List(ctx context.Context, params domain.ListParams) ([]domain.User, int, error) {
	return s.userRepo.List(ctx, params)
}

func (s *UserService) GetSchool(ctx context.Context, schoolID uuid.UUID) (*domain.School, error) {
	return s.schoolRepo.GetByID(ctx, schoolID)
}

func (s *UserService) UpdatePhoto(ctx context.Context, id uuid.UUID, photoURL string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.PhotoURL = &photoURL
	return s.userRepo.Update(ctx, user)
}
