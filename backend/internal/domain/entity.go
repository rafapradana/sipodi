package domain

import (
	"time"

	"github.com/google/uuid"
)

// Enums
type UserRole string

const (
	RoleSuperAdmin   UserRole = "super_admin"
	RoleAdminSekolah UserRole = "admin_sekolah"
	RoleGTK          UserRole = "gtk"
)

type Gender string

const (
	GenderMale   Gender = "L"
	GenderFemale Gender = "P"
)

type GTKType string

const (
	GTKTypeGuru          GTKType = "guru"
	GTKTypeTendik        GTKType = "tendik"
	GTKTypeKepalaSekolah GTKType = "kepala_sekolah"
)

type SchoolStatus string

const (
	SchoolStatusNegeri SchoolStatus = "negeri"
	SchoolStatusSwasta SchoolStatus = "swasta"
)

type TalentType string

const (
	TalentTypePesertaPelatihan TalentType = "peserta_pelatihan"
	TalentTypePembimbingLomba  TalentType = "pembimbing_lomba"
	TalentTypePesertaLomba     TalentType = "peserta_lomba"
	TalentTypeMinatBakat       TalentType = "minat_bakat"
)

type TalentStatus string

const (
	TalentStatusPending  TalentStatus = "pending"
	TalentStatusApproved TalentStatus = "approved"
	TalentStatusRejected TalentStatus = "rejected"
)

type CompetitionLevel string

const (
	LevelKota          CompetitionLevel = "kota"
	LevelProvinsi      CompetitionLevel = "provinsi"
	LevelNasional      CompetitionLevel = "nasional"
	LevelInternasional CompetitionLevel = "internasional"
)

type TalentField string

const (
	FieldAkademik     TalentField = "akademik"
	FieldInovasi      TalentField = "inovasi"
	FieldTeknologi    TalentField = "teknologi"
	FieldSosial       TalentField = "sosial"
	FieldOlahraga     TalentField = "olahraga"
	FieldSeni         TalentField = "seni"
	FieldKepemimpinan TalentField = "kepemimpinan"
)

type NotificationType string

const (
	NotificationTalentApproved NotificationType = "talent_approved"
	NotificationTalentRejected NotificationType = "talent_rejected"
)

// Entities
type School struct {
	ID           uuid.UUID    `json:"id"`
	Name         string       `json:"name"`
	NPSN         string       `json:"npsn"`
	Status       SchoolStatus `json:"status"`
	Address      string       `json:"address"`
	HeadMasterID *uuid.UUID   `json:"head_master_id,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         UserRole   `json:"role"`
	FullName     string     `json:"full_name"`
	PhotoURL     *string    `json:"photo_url,omitempty"`
	NUPTK        *string    `json:"nuptk,omitempty"`
	NIP          *string    `json:"nip,omitempty"`
	Gender       *Gender    `json:"gender,omitempty"`
	BirthDate    *time.Time `json:"birth_date,omitempty"`
	GTKType      *GTKType   `json:"gtk_type,omitempty"`
	Position     *string    `json:"position,omitempty"`
	SchoolID     *uuid.UUID `json:"school_id,omitempty"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Talent struct {
	ID              uuid.UUID    `json:"id"`
	UserID          uuid.UUID    `json:"user_id"`
	TalentType      TalentType   `json:"talent_type"`
	Status          TalentStatus `json:"status"`
	VerifiedBy      *uuid.UUID   `json:"verified_by,omitempty"`
	VerifiedAt      *time.Time   `json:"verified_at,omitempty"`
	RejectionReason *string      `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type TalentTraining struct {
	ID           uuid.UUID `json:"id"`
	TalentID     uuid.UUID `json:"talent_id"`
	ActivityName string    `json:"activity_name"`
	Organizer    string    `json:"organizer"`
	StartDate    time.Time `json:"start_date"`
	DurationDays int       `json:"duration_days"`
}

type TalentCompetitionMentor struct {
	ID              uuid.UUID        `json:"id"`
	TalentID        uuid.UUID        `json:"talent_id"`
	CompetitionName string           `json:"competition_name"`
	Level           CompetitionLevel `json:"level"`
	Organizer       string           `json:"organizer"`
	Field           TalentField      `json:"field"`
	Achievement     string           `json:"achievement"`
	CertificateURL  *string          `json:"certificate_url,omitempty"`
}

type TalentCompetitionParticipant struct {
	ID               uuid.UUID        `json:"id"`
	TalentID         uuid.UUID        `json:"talent_id"`
	CompetitionName  string           `json:"competition_name"`
	Level            CompetitionLevel `json:"level"`
	Organizer        string           `json:"organizer"`
	Field            TalentField      `json:"field"`
	StartDate        time.Time        `json:"start_date"`
	DurationDays     int              `json:"duration_days"`
	CompetitionField string           `json:"competition_field"`
	Achievement      string           `json:"achievement"`
	CertificateURL   *string          `json:"certificate_url,omitempty"`
}

type TalentInterest struct {
	ID             uuid.UUID `json:"id"`
	TalentID       uuid.UUID `json:"talent_id"`
	InterestName   string    `json:"interest_name"`
	Description    string    `json:"description"`
	CertificateURL *string   `json:"certificate_url,omitempty"`
}

type Notification struct {
	ID        uuid.UUID        `json:"id"`
	UserID    uuid.UUID        `json:"user_id"`
	TalentID  *uuid.UUID       `json:"talent_id,omitempty"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
}
