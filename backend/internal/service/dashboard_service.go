package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/repository"
)

type DashboardService struct {
	userRepo         *repository.UserRepository
	schoolRepo       *repository.SchoolRepository
	talentRepo       *repository.TalentRepository
	notificationRepo *repository.NotificationRepository
}

func NewDashboardService(
	userRepo *repository.UserRepository,
	schoolRepo *repository.SchoolRepository,
	talentRepo *repository.TalentRepository,
	notificationRepo *repository.NotificationRepository,
) *DashboardService {
	return &DashboardService{
		userRepo:         userRepo,
		schoolRepo:       schoolRepo,
		talentRepo:       talentRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *DashboardService) GetSuperAdminSummary(ctx context.Context) (*domain.DashboardSummary, error) {
	totalSchools, _ := s.schoolRepo.Count(ctx)
	totalGTK, _ := s.userRepo.CountByRole(ctx, domain.RoleGTK)
	totalAdminSekolah, _ := s.userRepo.CountByRole(ctx, domain.RoleAdminSekolah)
	gtkByType, _ := s.userRepo.CountByGTKType(ctx)
	totalTalents, _ := s.talentRepo.Count(ctx)
	talentsByStatus, _ := s.talentRepo.CountByStatus(ctx)
	talentsByType, _ := s.talentRepo.CountByType(ctx)

	return &domain.DashboardSummary{
		TotalSchools:      totalSchools,
		TotalUsers:        totalGTK + totalAdminSekolah + 1, // +1 for super admin
		TotalGTK:          totalGTK,
		TotalAdminSekolah: totalAdminSekolah,
		GTKByType:         gtkByType,
		TotalTalents:      totalTalents,
		TalentsByStatus:   talentsByStatus,
		TalentsByType:     talentsByType,
	}, nil
}

func (s *DashboardService) GetAdminSekolahSummary(ctx context.Context, schoolID uuid.UUID) (*domain.DashboardSummary, error) {
	school, _ := s.schoolRepo.GetByID(ctx, schoolID)
	totalGTK, _ := s.userRepo.CountBySchool(ctx, schoolID)

	guruCount, _ := s.userRepo.CountBySchoolAndType(ctx, schoolID, domain.GTKTypeGuru)
	tendikCount, _ := s.userRepo.CountBySchoolAndType(ctx, schoolID, domain.GTKTypeTendik)
	kepalaCount, _ := s.userRepo.CountBySchoolAndType(ctx, schoolID, domain.GTKTypeKepalaSekolah)

	talentsByStatus, _ := s.talentRepo.CountBySchoolID(ctx, schoolID)
	pendingVerifications, _ := s.talentRepo.CountPendingBySchoolID(ctx, schoolID)

	var schoolRef *domain.SchoolRef
	if school != nil {
		schoolRef = &domain.SchoolRef{
			ID:   school.ID,
			Name: school.Name,
		}
	}

	return &domain.DashboardSummary{
		School:   schoolRef,
		TotalGTK: totalGTK,
		GTKByType: map[string]int{
			"guru":           guruCount,
			"tendik":         tendikCount,
			"kepala_sekolah": kepalaCount,
		},
		TotalTalents:         talentsByStatus["total"],
		TalentsByStatus:      talentsByStatus,
		PendingVerifications: pendingVerifications,
	}, nil
}

func (s *DashboardService) GetGTKSummary(ctx context.Context, userID uuid.UUID) (*domain.DashboardSummary, error) {
	myTalents, _ := s.talentRepo.CountByUserID(ctx, userID)
	unreadNotifications, _ := s.notificationRepo.CountUnread(ctx, userID)

	return &domain.DashboardSummary{
		MyTalents:           myTalents,
		UnreadNotifications: unreadNotifications,
	}, nil
}

func (s *DashboardService) GetSchoolsStatistics(ctx context.Context, params domain.ListParams) ([]domain.SchoolStatistics, int, error) {
	return s.schoolRepo.GetStatistics(ctx, params)
}

func (s *DashboardService) GetTalentsStatistics(ctx context.Context, groupBy string, schoolID *string, dateFrom, dateTo string) (map[string]interface{}, error) {
	return s.talentRepo.GetStatistics(ctx, groupBy, schoolID, dateFrom, dateTo)
}
