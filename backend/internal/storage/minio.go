package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/Helltale/beer-mania/backend/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOStorage implements Storage interface using MinIO client
type MinIOStorage struct {
	client *minio.Client
	cfg    *config.MinIOConfig
}

// NewMinIOStorage creates a new MinIO storage client
func NewMinIOStorage(cfg *config.MinIOConfig) (*MinIOStorage, error) {
	// Initialize MinIO client
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	storage := &MinIOStorage{
		client: client,
		cfg:    cfg,
	}

	// Ensure buckets exist
	ctx := context.Background()
	if err := storage.EnsureBucketExists(ctx, cfg.BucketUploads); err != nil {
		return nil, fmt.Errorf("failed to ensure uploads bucket exists: %w", err)
	}

	if err := storage.EnsureBucketExists(ctx, cfg.BucketProcessed); err != nil {
		return nil, fmt.Errorf("failed to ensure processed bucket exists: %w", err)
	}

	return storage, nil
}

// EnsureBucketExists ensures that a bucket exists, creates it if it doesn't
func (s *MinIOStorage) EnsureBucketExists(ctx context.Context, bucketName string) error {
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		if err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
	}

	return nil
}

// UploadFile uploads a file to MinIO storage and returns the object URL
func (s *MinIOStorage) UploadFile(ctx context.Context, bucket string, objectName string, file io.Reader, size int64, contentType string) (string, error) {
	// Upload file
	_, err := s.client.PutObject(ctx, bucket, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate URL
	url, err := s.GetFileURL(ctx, bucket, objectName)
	if err != nil {
		return "", fmt.Errorf("failed to generate file URL: %w", err)
	}

	return url, nil
}

// GetFileURL returns the URL to access a file in MinIO storage
// Returns a presigned URL valid for configured expiration time
func (s *MinIOStorage) GetFileURL(ctx context.Context, bucket string, objectName string) (string, error) {
	// Generate a presigned URL with expiration time from config
	presignedURL, err := s.client.PresignedGetObject(ctx, bucket, objectName, s.cfg.PresignedURLExpiration(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

// DeleteFile deletes a file from MinIO storage
func (s *MinIOStorage) DeleteFile(ctx context.Context, bucket string, objectName string) error {
	if err := s.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
