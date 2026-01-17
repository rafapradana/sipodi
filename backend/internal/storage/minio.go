package storage

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sipodi/backend/internal/config"
)

type MinIOStorage struct {
	client       *minio.Client
	signerClient *minio.Client
	bucket       string
	publicURL    string
}

func NewMinIOStorage(cfg config.MinIOConfig) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Create a separate client for signing URLs using the Public URL (browser accessible)
	// This ensures the Host header in the signature matches what the browser sends.
	// Presigning is an offline operation, so this client doesn't need to be able to connect mostly,
	// but it needs the correct endpoint config.
	var signerClient *minio.Client
	if cfg.PublicURL != "" {
		u, err := url.Parse(cfg.PublicURL)
		if err != nil {
			return nil, fmt.Errorf("invalid public url: %w", err)
		}
		secure := u.Scheme == "https"
		endpoint := u.Host

		signerClient, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: secure,
			Region: "us-east-1", // Validate that this prevents lookup
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create signer client: %w", err)
		}
	} else {
		// Fallback to internal client if no public URL configured
		signerClient = client
	}

	return &MinIOStorage{
		client:       client,
		signerClient: signerClient,
		bucket:       cfg.Bucket,
		publicURL:    cfg.PublicURL,
	}, nil
}

func (s *MinIOStorage) GeneratePresignedURL(ctx context.Context, objectName string, contentType string, expiry time.Duration) (string, error) {
	// Use signerClient to generate URL with public endpoint logic
	presignedURL, err := s.signerClient.PresignedPutObject(ctx, s.bucket, objectName, expiry)
	if err != nil {
		fmt.Printf("ERROR generating presigned URL: %v\n", err)
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedURL.String(), nil
}

func (s *MinIOStorage) GetObjectURL(objectName string) string {
	return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, objectName)
}

func (s *MinIOStorage) ObjectExists(ctx context.Context, objectName string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *MinIOStorage) GetObjectInfo(ctx context.Context, objectName string) (minio.ObjectInfo, error) {
	return s.client.StatObject(ctx, s.bucket, objectName, minio.StatObjectOptions{})
}

func (s *MinIOStorage) DeleteObject(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}

func (s *MinIOStorage) GenerateObjectName(uploadType string, filename string) string {
	now := time.Now()
	ext := path.Ext(filename)
	id := uuid.New().String()
	return fmt.Sprintf("%s/%d/%02d/%s%s", uploadType, now.Year(), now.Month(), id, ext)
}

func (s *MinIOStorage) GetPresignedGetURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	// Use signerClient here too
	presignedURL, err := s.signerClient.PresignedGetObject(ctx, s.bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
