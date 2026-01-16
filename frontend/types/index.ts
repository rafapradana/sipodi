// =================================
// SIPODI TypeScript Types
// Based on backend/internal/domain
// =================================

// Enums
export type UserRole = 'super_admin' | 'admin_sekolah' | 'gtk';
export type Gender = 'L' | 'P';
export type GTKType = 'guru' | 'tendik' | 'kepala_sekolah';
export type SchoolStatus = 'negeri' | 'swasta';
export type TalentType = 'peserta_pelatihan' | 'pembimbing_lomba' | 'peserta_lomba' | 'minat_bakat';
export type TalentStatus = 'pending' | 'approved' | 'rejected';
export type CompetitionLevel = 'kota' | 'provinsi' | 'nasional' | 'internasional';
export type TalentField = 'akademik' | 'inovasi' | 'teknologi' | 'sosial' | 'olahraga' | 'seni' | 'kepemimpinan';
export type NotificationType = 'talent_approved' | 'talent_rejected';

// User Roles Display
export const USER_ROLE_LABELS: Record<UserRole, string> = {
  super_admin: 'Super Admin',
  admin_sekolah: 'Admin Sekolah',
  gtk: 'GTK',
};

export const GTK_TYPE_LABELS: Record<GTKType, string> = {
  guru: 'Guru',
  tendik: 'Tenaga Kependidikan',
  kepala_sekolah: 'Kepala Sekolah',
};

export const SCHOOL_STATUS_LABELS: Record<SchoolStatus, string> = {
  negeri: 'Negeri',
  swasta: 'Swasta',
};

export const TALENT_TYPE_LABELS: Record<TalentType, string> = {
  peserta_pelatihan: 'Peserta Pelatihan',
  pembimbing_lomba: 'Pembimbing Lomba',
  peserta_lomba: 'Peserta Lomba',
  minat_bakat: 'Minat/Bakat',
};

export const TALENT_STATUS_LABELS: Record<TalentStatus, string> = {
  pending: 'Pending',
  approved: 'Disetujui',
  rejected: 'Ditolak',
};

export const COMPETITION_LEVEL_LABELS: Record<CompetitionLevel, string> = {
  kota: 'Kota',
  provinsi: 'Provinsi',
  nasional: 'Nasional',
  internasional: 'Internasional',
};

export const TALENT_FIELD_LABELS: Record<TalentField, string> = {
  akademik: 'Akademik',
  inovasi: 'Inovasi',
  teknologi: 'Teknologi',
  sosial: 'Sosial',
  olahraga: 'Olahraga',
  seni: 'Seni',
  kepemimpinan: 'Kepemimpinan',
};

// Entities
export interface SchoolRef {
  id: string;
  name: string;
  npsn?: string;
}

export interface UserRef {
  id: string;
  full_name: string;
  nip?: string;
}

export interface School {
  id: string;
  name: string;
  npsn: string;
  status: SchoolStatus;
  address: string;
  head_master?: UserRef;
  gtk_count: number;
  created_at: string;
  updated_at: string;
}

export interface SchoolDetail extends School {
  guru_count: number;
  tendik_count: number;
  kepala_sekolah_count: number;
}

export interface User {
  id: string;
  email: string;
  role: UserRole;
  full_name: string;
  photo_url?: string;
  nuptk?: string;
  nip?: string;
  gender?: Gender;
  birth_date?: string;
  gtk_type?: GTKType;
  position?: string;
  is_active: boolean;
  school?: SchoolRef;
  created_at: string;
  updated_at: string;
}

export interface UserListItem {
  id: string;
  email: string;
  role: UserRole;
  full_name: string;
  nuptk?: string;
  nip?: string;
  gtk_type?: GTKType;
  position?: string;
  is_active: boolean;
  school?: SchoolRef;
  created_at: string;
}

export interface TalentUser {
  id: string;
  full_name: string;
  school_name?: string;
}

export interface Talent {
  id: string;
  user?: TalentUser;
  talent_type: TalentType;
  status: TalentStatus;
  detail: TrainingDetail | MentorDetail | ParticipantDetail | InterestDetail;
  certificate_url?: string;
  verified_by?: UserRef;
  verified_at?: string;
  rejection_reason?: string;
  created_at: string;
  updated_at: string;
}

export interface TalentListItem {
  id: string;
  user?: TalentUser;
  talent_type: TalentType;
  status: TalentStatus;
  detail: TrainingDetail | MentorDetail | ParticipantDetail | InterestDetail;
  created_at: string;
  updated_at: string;
}

// Talent Detail Types
export interface TrainingDetail {
  activity_name: string;
  organizer: string;
  start_date: string;
  duration_days: number;
}

export interface MentorDetail {
  competition_name: string;
  level: CompetitionLevel;
  organizer: string;
  field: TalentField;
  achievement: string;
}

export interface ParticipantDetail {
  competition_name: string;
  level: CompetitionLevel;
  organizer: string;
  field: TalentField;
  start_date: string;
  duration_days: number;
  competition_field: string;
  achievement: string;
}

export interface InterestDetail {
  interest_name: string;
  description: string;
}

export interface Notification {
  id: string;
  type: NotificationType;
  message: string;
  talent_id?: string;
  is_read: boolean;
  created_at: string;
}

// Request DTOs
export interface LoginRequest {
  email: string;
  password: string;
}

export interface CreateSchoolRequest {
  name: string;
  npsn: string;
  status: SchoolStatus;
  address: string;
}

export interface UpdateSchoolRequest {
  name?: string;
  npsn?: string;
  status?: SchoolStatus;
  address?: string;
  head_master_id?: string;
}

export interface CreateUserRequest {
  email: string;
  password: string;
  role: UserRole;
  full_name: string;
  nuptk?: string;
  nip?: string;
  gender?: Gender;
  birth_date?: string;
  gtk_type?: GTKType;
  position?: string;
  school_id?: string;
}

export interface UpdateUserRequest {
  full_name?: string;
  nuptk?: string;
  nip?: string;
  gender?: Gender;
  birth_date?: string;
  gtk_type?: GTKType;
  position?: string;
  school_id?: string;
}

export interface UpdateTalentRequest {
  talent_type?: TalentType;
  detail?: TrainingDetail | MentorDetail | ParticipantDetail | InterestDetail;
  upload_id?: string;
}

export interface RejectTalentRequest {
  rejection_reason: string;
}

export interface BatchApproveRequest {
  ids: string[];
}

export interface BatchRejectRequest {
  ids: string[];
  rejection_reason: string;
}

export interface PresignResponse {
  upload_id: string;
  presigned_url: string;
  method: string;
  expires_in: number;
  max_size: number;
  allowed_types: string[];
}

// Response DTOs
export interface LoginResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}

export interface RefreshResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
}

export interface PaginationMeta {
  current_page: number;
  per_page: number;
  total_pages: number;
  total_count: number;
}

export interface ListResponse<T> {
  data: T[];
  meta: PaginationMeta;
}

export interface DataResponse<T> {
  data: T;
  message?: string;
}

export interface ApiError {
  error: {
    code: string;
    message: string;
    details?: { field: string; message: string }[];
  };
}

export interface BatchResult {
  approved_count?: number;
  rejected_count?: number;
  failed_count: number;
  failed_ids: { id: string; reason: string }[];
}

// Dashboard DTOs
export interface DashboardSummary {
  // Super Admin
  total_schools?: number;
  total_users?: number;
  total_gtk?: number;
  total_admin_sekolah?: number;
  gtk_by_type?: Record<string, number>;
  total_talents?: number;
  talents_by_status?: Record<string, number>;
  talents_by_type?: Record<string, number>;
  recent_talents?: TalentListItem[];

  // Admin Sekolah
  school?: SchoolRef;
  pending_verifications?: number;

  // GTK
  my_talents?: Record<string, number>;
  unread_notifications?: number;
}

export interface SchoolStatistics {
  id: string;
  name: string;
  npsn: string;
  status: SchoolStatus;
  gtk_count: number;
  talent_count: number;
  pending_count: number;
  approved_count: number;
}
