-- ============================================
-- SIPODI Database Schema
-- PostgreSQL - BCNF Normalized
-- ============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- ENUM TYPES
-- ============================================

CREATE TYPE user_role AS ENUM ('super_admin', 'admin_sekolah', 'gtk');
CREATE TYPE gender AS ENUM ('L', 'P');
CREATE TYPE gtk_type AS ENUM ('guru', 'tendik', 'kepala_sekolah');
CREATE TYPE school_status AS ENUM ('negeri', 'swasta');
CREATE TYPE talent_type AS ENUM ('peserta_pelatihan', 'pembimbing_lomba', 'peserta_lomba', 'minat_bakat');
CREATE TYPE talent_status AS ENUM ('pending', 'approved', 'rejected');
CREATE TYPE competition_level AS ENUM ('kota', 'provinsi', 'nasional', 'internasional');
CREATE TYPE talent_field AS ENUM ('akademik', 'inovasi', 'teknologi', 'sosial', 'olahraga', 'seni', 'kepemimpinan');
CREATE TYPE notification_type AS ENUM ('talent_approved', 'talent_rejected');

-- ============================================
-- TABLES
-- ============================================

-- Schools table
CREATE TABLE schools (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    npsn VARCHAR(20) UNIQUE NOT NULL,
    status school_status NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Users table (GTK, Admin Sekolah, Super Admin)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'gtk',
    full_name VARCHAR(255) NOT NULL,
    photo_url VARCHAR(500),
    nuptk VARCHAR(20) UNIQUE,
    nip VARCHAR(20) UNIQUE,
    gender gender,
    birth_date DATE,
    gtk_type gtk_type,
    position VARCHAR(255),
    school_id UUID REFERENCES schools(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Update schools to reference head master (after users table exists)
ALTER TABLE schools 
ADD COLUMN head_master_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- Refresh tokens table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- TALENT TABLES (Normalized by type)
-- ============================================

-- Base talents table
CREATE TABLE talents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    talent_type talent_type NOT NULL,
    status talent_status DEFAULT 'pending',
    verified_by UUID REFERENCES users(id) ON DELETE SET NULL,
    verified_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Peserta Pelatihan
CREATE TABLE talent_trainings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    talent_id UUID UNIQUE NOT NULL REFERENCES talents(id) ON DELETE CASCADE,
    activity_name VARCHAR(255) NOT NULL,
    organizer VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    duration_days INTEGER NOT NULL CHECK (duration_days > 0)
);

-- Pembimbing Lomba
CREATE TABLE talent_competition_mentors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    talent_id UUID UNIQUE NOT NULL REFERENCES talents(id) ON DELETE CASCADE,
    competition_name VARCHAR(255) NOT NULL,
    level competition_level NOT NULL,
    organizer VARCHAR(255) NOT NULL,
    field talent_field NOT NULL,
    achievement TEXT NOT NULL,
    certificate_url VARCHAR(500)
);

-- Peserta Lomba
CREATE TABLE talent_competition_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    talent_id UUID UNIQUE NOT NULL REFERENCES talents(id) ON DELETE CASCADE,
    competition_name VARCHAR(255) NOT NULL,
    level competition_level NOT NULL,
    organizer VARCHAR(255) NOT NULL,
    field talent_field NOT NULL,
    start_date DATE NOT NULL,
    duration_days INTEGER NOT NULL CHECK (duration_days > 0),
    competition_field VARCHAR(255) NOT NULL,
    achievement TEXT NOT NULL,
    certificate_url VARCHAR(500)
);

-- Minat/Bakat
CREATE TABLE talent_interests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    talent_id UUID UNIQUE NOT NULL REFERENCES talents(id) ON DELETE CASCADE,
    interest_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    certificate_url VARCHAR(500)
);

-- ============================================
-- NOTIFICATIONS
-- ============================================

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    talent_id UUID REFERENCES talents(id) ON DELETE CASCADE,
    type notification_type NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- INDEXES
-- ============================================

-- Users indexes
CREATE INDEX idx_users_school_id ON users(school_id);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_gtk_type ON users(gtk_type);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_nuptk ON users(nuptk);
CREATE INDEX idx_users_nip ON users(nip);

-- Schools indexes
CREATE INDEX idx_schools_npsn ON schools(npsn);
CREATE INDEX idx_schools_status ON schools(status);

-- Talents indexes
CREATE INDEX idx_talents_user_id ON talents(user_id);
CREATE INDEX idx_talents_status ON talents(status);
CREATE INDEX idx_talents_type ON talents(talent_type);
CREATE INDEX idx_talents_user_status ON talents(user_id, status);

-- Refresh tokens indexes
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- Notifications indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(user_id, is_read);

-- ============================================
-- FUNCTIONS & TRIGGERS
-- ============================================

-- Auto update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to tables with updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_schools_updated_at
    BEFORE UPDATE ON schools
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_talents_updated_at
    BEFORE UPDATE ON talents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to create notification when talent status changes
CREATE OR REPLACE FUNCTION notify_talent_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status = 'pending' AND NEW.status IN ('approved', 'rejected') THEN
        INSERT INTO notifications (user_id, talent_id, type, message)
        VALUES (
            NEW.user_id,
            NEW.id,
            CASE WHEN NEW.status = 'approved' THEN 'talent_approved'::notification_type 
                 ELSE 'talent_rejected'::notification_type END,
            CASE WHEN NEW.status = 'approved' 
                 THEN 'Talenta Anda telah disetujui'
                 ELSE 'Talenta Anda ditolak. Alasan: ' || COALESCE(NEW.rejection_reason, '-') END
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_talent_status_notification
    AFTER UPDATE ON talents
    FOR EACH ROW
    WHEN (OLD.status IS DISTINCT FROM NEW.status)
    EXECUTE FUNCTION notify_talent_status_change();

-- Function to reset talent status to pending on update
CREATE OR REPLACE FUNCTION reset_talent_status_on_update()
RETURNS TRIGGER AS $$
BEGIN
    -- Reset status to pending when talent detail is updated
    UPDATE talents 
    SET status = 'pending', 
        verified_by = NULL, 
        verified_at = NULL,
        rejection_reason = NULL
    WHERE id = NEW.talent_id 
    AND status != 'pending';
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply reset trigger to all talent detail tables
CREATE TRIGGER reset_training_status
    AFTER UPDATE ON talent_trainings
    FOR EACH ROW
    EXECUTE FUNCTION reset_talent_status_on_update();

CREATE TRIGGER reset_mentor_status
    AFTER UPDATE ON talent_competition_mentors
    FOR EACH ROW
    EXECUTE FUNCTION reset_talent_status_on_update();

CREATE TRIGGER reset_participant_status
    AFTER UPDATE ON talent_competition_participants
    FOR EACH ROW
    EXECUTE FUNCTION reset_talent_status_on_update();

CREATE TRIGGER reset_interest_status
    AFTER UPDATE ON talent_interests
    FOR EACH ROW
    EXECUTE FUNCTION reset_talent_status_on_update();

-- ============================================
-- SEED DATA (Super Admin)
-- ============================================

-- Default super admin (password: admin123 - hashed with bcrypt)
-- Note: Change this password in production!
INSERT INTO users (email, password_hash, role, full_name)
VALUES (
    'superadmin@sipodi.go.id',
    '$2a$10$dXG3Jih7K.UvR5mYXAaeoOb4TAqTqp9JR170UuLVaYkm79xsxsyHi',
    'super_admin',
    'Super Administrator'
);
