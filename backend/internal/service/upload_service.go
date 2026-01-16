package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/storage"
)

var (
	ErrUploadNotFound  = errors.New("upload not found")
	ErrFileNotUploaded = errors.New("file not uploaded")
	ErrInvalidFileType = errors.New("invalid file type")
)

type UploadInfo struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ObjectName  string
	Filename    string
	ContentType string
	UploadType  string
	ExpiresAt   time.Time
}

type UploadService struct {
	storage *storage.MinIOStorage
	uploads map[uuid.UUID]*UploadInfo
	mu      sync.RWMutex
}

func NewUploadService(storage *storage.MinIOStorage) *UploadService {
	return &UploadService{
		storage: storage,
		uploads: make(map[uuid.UUID]*UploadInfo),
	}
}

var uploadTypeConfig = map[string]struct {
	MaxSize      int64
	AllowedTypes []string
}{
	"profile_photo": {
		MaxSize:      2 * 1024 * 1024, // 2MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp"},
	},
	"talent_certificate": {
		MaxSize:      10 * 1024 * 1024, // 10MB
		AllowedTypes: []string{"application/pdf", "image/jpeg", "image/png"},
	},
}

func (s *UploadService) GeneratePresignedURL(ctx context.Context, userID uuid.UUID, req domain.PresignRequest) (*domain.PresignResponse, error) {
	config, ok := uploadTypeConfig[req.UploadType]
	if !ok {
		return nil, ErrInvalidFileType
	}

	// Validate content type
	validType := false
	for _, t := range config.AllowedTypes {
		if t == req.ContentType {
			validType = true
			break
		}
	}
	if !validType {
		return nil, ErrInvalidFileType
	}

	uploadID := uuid.New()
	objectName := s.storage.GenerateObjectName(req.UploadType, req.Filename)
	expiry := time.Hour

	presignedURL, err := s.storage.GeneratePresignedURL(ctx, objectName, req.ContentType, expiry)
	if err != nil {
		return nil, err
	}

	// Store upload info
	s.mu.Lock()
	s.uploads[uploadID] = &UploadInfo{
		ID:          uploadID,
		UserID:      userID,
		ObjectName:  objectName,
		Filename:    req.Filename,
		ContentType: req.ContentType,
		UploadType:  req.UploadType,
		ExpiresAt:   time.Now().Add(expiry),
	}
	s.mu.Unlock()

	return &domain.PresignResponse{
		UploadID:     uploadID,
		PresignedURL: presignedURL,
		Method:       "PUT",
		ExpiresIn:    int(expiry.Seconds()),
		MaxSize:      config.MaxSize,
		AllowedTypes: config.AllowedTypes,
	}, nil
}

func (s *UploadService) ConfirmUpload(ctx context.Context, uploadID uuid.UUID, userID uuid.UUID) (*domain.ConfirmUploadResponse, error) {
	s.mu.RLock()
	info, ok := s.uploads[uploadID]
	s.mu.RUnlock()

	if !ok || info.UserID != userID {
		return nil, ErrUploadNotFound
	}

	if time.Now().After(info.ExpiresAt) {
		s.mu.Lock()
		delete(s.uploads, uploadID)
		s.mu.Unlock()
		return nil, ErrUploadNotFound
	}

	// Check if file exists in MinIO
	exists, err := s.storage.ObjectExists(ctx, info.ObjectName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrFileNotUploaded
	}

	// Get file info
	objInfo, err := s.storage.GetObjectInfo(ctx, info.ObjectName)
	if err != nil {
		return nil, err
	}

	// Remove from pending uploads
	s.mu.Lock()
	delete(s.uploads, uploadID)
	s.mu.Unlock()

	return &domain.ConfirmUploadResponse{
		UploadID:    uploadID,
		FileURL:     s.storage.GetObjectURL(info.ObjectName),
		Filename:    info.Filename,
		FileSize:    objInfo.Size,
		ContentType: info.ContentType,
	}, nil
}

func (s *UploadService) CancelUpload(ctx context.Context, uploadID uuid.UUID, userID uuid.UUID) error {
	s.mu.RLock()
	info, ok := s.uploads[uploadID]
	s.mu.RUnlock()

	if !ok || info.UserID != userID {
		return ErrUploadNotFound
	}

	// Try to delete from MinIO if uploaded
	s.storage.DeleteObject(ctx, info.ObjectName)

	s.mu.Lock()
	delete(s.uploads, uploadID)
	s.mu.Unlock()

	return nil
}

func (s *UploadService) GetUploadInfo(uploadID uuid.UUID) (*UploadInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.uploads[uploadID]
	return info, ok
}

func (s *UploadService) GetFileURL(objectName string) string {
	return s.storage.GetObjectURL(objectName)
}

// Cleanup expired uploads periodically
func (s *UploadService) CleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, info := range s.uploads {
		if now.After(info.ExpiresAt) {
			delete(s.uploads, id)
		}
	}
}
