package domain

import (
	"time"

	"github.com/google/uuid"
)

// Auth DTOs
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int          `json:"expires_in"`
	User        UserResponse `json:"user"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// User DTOs
type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Role      UserRole   `json:"role"`
	FullName  string     `json:"full_name"`
	PhotoURL  *string    `json:"photo_url,omitempty"`
	NUPTK     *string    `json:"nuptk,omitempty"`
	NIP       *string    `json:"nip,omitempty"`
	Gender    *Gender    `json:"gender,omitempty"`
	BirthDate *string    `json:"birth_date,omitempty"`
	GTKType   *GTKType   `json:"gtk_type,omitempty"`
	Position  *string    `json:"position,omitempty"`
	IsActive  bool       `json:"is_active"`
	School    *SchoolRef `json:"school,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type UserListResponse struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Role      UserRole   `json:"role"`
	FullName  string     `json:"full_name"`
	NUPTK     *string    `json:"nuptk,omitempty"`
	NIP       *string    `json:"nip,omitempty"`
	GTKType   *GTKType   `json:"gtk_type,omitempty"`
	Position  *string    `json:"position,omitempty"`
	IsActive  bool       `json:"is_active"`
	School    *SchoolRef `json:"school,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateUserRequest struct {
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      UserRole   `json:"role"`
	FullName  string     `json:"full_name"`
	NUPTK     *string    `json:"nuptk,omitempty"`
	NIP       *string    `json:"nip,omitempty"`
	Gender    *Gender    `json:"gender,omitempty"`
	BirthDate *string    `json:"birth_date,omitempty"`
	GTKType   *GTKType   `json:"gtk_type,omitempty"`
	Position  *string    `json:"position,omitempty"`
	SchoolID  *uuid.UUID `json:"school_id,omitempty"`
}

type UpdateUserRequest struct {
	FullName  *string    `json:"full_name,omitempty"`
	NUPTK     *string    `json:"nuptk,omitempty"`
	NIP       *string    `json:"nip,omitempty"`
	Gender    *Gender    `json:"gender,omitempty"`
	BirthDate *string    `json:"birth_date,omitempty"`
	GTKType   *GTKType   `json:"gtk_type,omitempty"`
	Position  *string    `json:"position,omitempty"`
	SchoolID  *uuid.UUID `json:"school_id,omitempty"`
}

type UpdateProfileRequest struct {
	FullName  *string `json:"full_name,omitempty"`
	Gender    *Gender `json:"gender,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"`
	Position  *string `json:"position,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword         string `json:"current_password"`
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

// School DTOs
type SchoolRef struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	NPSN string    `json:"npsn,omitempty"`
}

type SchoolResponse struct {
	ID         uuid.UUID    `json:"id"`
	Name       string       `json:"name"`
	NPSN       string       `json:"npsn"`
	Status     SchoolStatus `json:"status"`
	Address    string       `json:"address"`
	HeadMaster *UserRef     `json:"head_master,omitempty"`
	GTKCount   int          `json:"gtk_count"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type SchoolDetailResponse struct {
	ID                 uuid.UUID    `json:"id"`
	Name               string       `json:"name"`
	NPSN               string       `json:"npsn"`
	Status             SchoolStatus `json:"status"`
	Address            string       `json:"address"`
	HeadMaster         *UserRef     `json:"head_master,omitempty"`
	GTKCount           int          `json:"gtk_count"`
	GuruCount          int          `json:"guru_count"`
	TendikCount        int          `json:"tendik_count"`
	KepalaSekolahCount int          `json:"kepala_sekolah_count"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
}

type CreateSchoolRequest struct {
	Name    string       `json:"name"`
	NPSN    string       `json:"npsn"`
	Status  SchoolStatus `json:"status"`
	Address string       `json:"address"`
}

type UpdateSchoolRequest struct {
	Name         *string       `json:"name,omitempty"`
	NPSN         *string       `json:"npsn,omitempty"`
	Status       *SchoolStatus `json:"status,omitempty"`
	Address      *string       `json:"address,omitempty"`
	HeadMasterID *uuid.UUID    `json:"head_master_id,omitempty"`
}

type UserRef struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	NIP      *string   `json:"nip,omitempty"`
}

// Talent DTOs
type TalentResponse struct {
	ID              uuid.UUID    `json:"id"`
	User            *UserRef     `json:"user,omitempty"`
	TalentType      TalentType   `json:"talent_type"`
	Status          TalentStatus `json:"status"`
	Detail          interface{}  `json:"detail"`
	CertificateURL  *string      `json:"certificate_url,omitempty"`
	VerifiedBy      *UserRef     `json:"verified_by,omitempty"`
	VerifiedAt      *time.Time   `json:"verified_at,omitempty"`
	RejectionReason *string      `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type TalentListResponse struct {
	ID         uuid.UUID    `json:"id"`
	User       *TalentUser  `json:"user,omitempty"`
	TalentType TalentType   `json:"talent_type"`
	Status     TalentStatus `json:"status"`
	Detail     interface{}  `json:"detail"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type TalentUser struct {
	ID         uuid.UUID `json:"id"`
	FullName   string    `json:"full_name"`
	SchoolName string    `json:"school_name,omitempty"`
}

type CreateTalentRequest struct {
	TalentType TalentType  `json:"talent_type"`
	Detail     interface{} `json:"detail"`
	UploadID   *uuid.UUID  `json:"upload_id,omitempty"`
}

type UpdateTalentRequest struct {
	Detail   interface{} `json:"detail"`
	UploadID *uuid.UUID  `json:"upload_id,omitempty"`
}

// Talent Detail DTOs
type TrainingDetail struct {
	ActivityName string `json:"activity_name"`
	Organizer    string `json:"organizer"`
	StartDate    string `json:"start_date"`
	DurationDays int    `json:"duration_days"`
}

type MentorDetail struct {
	CompetitionName string           `json:"competition_name"`
	Level           CompetitionLevel `json:"level"`
	Organizer       string           `json:"organizer"`
	Field           TalentField      `json:"field"`
	Achievement     string           `json:"achievement"`
}

type ParticipantDetail struct {
	CompetitionName  string           `json:"competition_name"`
	Level            CompetitionLevel `json:"level"`
	Organizer        string           `json:"organizer"`
	Field            TalentField      `json:"field"`
	StartDate        string           `json:"start_date"`
	DurationDays     int              `json:"duration_days"`
	CompetitionField string           `json:"competition_field"`
	Achievement      string           `json:"achievement"`
}

type InterestDetail struct {
	InterestName string `json:"interest_name"`
	Description  string `json:"description"`
}

// Verification DTOs
type RejectTalentRequest struct {
	RejectionReason string `json:"rejection_reason"`
}

type BatchApproveRequest struct {
	IDs []uuid.UUID `json:"ids"`
}

type BatchRejectRequest struct {
	IDs             []uuid.UUID `json:"ids"`
	RejectionReason string      `json:"rejection_reason"`
}

type BatchResult struct {
	SuccessCount int          `json:"approved_count,omitempty"`
	FailedCount  int          `json:"failed_count"`
	FailedIDs    []FailedItem `json:"failed_ids"`
}

type FailedItem struct {
	ID     uuid.UUID `json:"id"`
	Reason string    `json:"reason"`
}

// Notification DTOs
type NotificationResponse struct {
	ID        uuid.UUID        `json:"id"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	TalentID  *uuid.UUID       `json:"talent_id,omitempty"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
}

// Upload DTOs
type PresignRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	UploadType  string `json:"upload_type"`
}

type PresignResponse struct {
	UploadID     uuid.UUID `json:"upload_id"`
	PresignedURL string    `json:"presigned_url"`
	Method       string    `json:"method"`
	ExpiresIn    int       `json:"expires_in"`
	MaxSize      int64     `json:"max_size"`
	AllowedTypes []string  `json:"allowed_types"`
}

type ConfirmUploadResponse struct {
	UploadID    uuid.UUID `json:"upload_id"`
	FileURL     string    `json:"file_url"`
	Filename    string    `json:"filename"`
	FileSize    int64     `json:"file_size"`
	ContentType string    `json:"content_type"`
}

// Dashboard DTOs
type DashboardSummary struct {
	TotalSchools         int                  `json:"total_schools,omitempty"`
	TotalUsers           int                  `json:"total_users,omitempty"`
	TotalGTK             int                  `json:"total_gtk,omitempty"`
	TotalAdminSekolah    int                  `json:"total_admin_sekolah,omitempty"`
	GTKByType            map[string]int       `json:"gtk_by_type,omitempty"`
	TotalTalents         int                  `json:"total_talents,omitempty"`
	TalentsByStatus      map[string]int       `json:"talents_by_status,omitempty"`
	TalentsByType        map[string]int       `json:"talents_by_type,omitempty"`
	School               *SchoolRef           `json:"school,omitempty"`
	MyTalents            map[string]int       `json:"my_talents,omitempty"`
	PendingVerifications int                  `json:"pending_verifications,omitempty"`
	UnreadNotifications  int                  `json:"unread_notifications,omitempty"`
	RecentTalents        []TalentListResponse `json:"recent_talents,omitempty"`
}

// Pagination
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
}

type ListParams struct {
	Page    int
	Limit   int
	Search  string
	Sort    string
	Filters map[string]string
}

// Response wrappers
type DataResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type ListResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// School Statistics DTO
type SchoolStatistics struct {
	ID            uuid.UUID    `json:"id"`
	Name          string       `json:"name"`
	NPSN          string       `json:"npsn"`
	Status        SchoolStatus `json:"status"`
	GTKCount      int          `json:"gtk_count"`
	TalentCount   int          `json:"talent_count"`
	PendingCount  int          `json:"pending_count"`
	ApprovedCount int          `json:"approved_count"`
}
