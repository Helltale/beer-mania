package storage

import (
	"context"
	"io"
)

// Storage defines the interface for object storage operations
type Storage interface {
	// UploadFile uploads a file to the storage and returns the object URL
	UploadFile(ctx context.Context, bucket string, objectName string, file io.Reader, size int64, contentType string) (string, error)

	// GetFileURL returns the URL to access a file in the storage
	GetFileURL(ctx context.Context, bucket string, objectName string) (string, error)

	// DeleteFile deletes a file from the storage
	DeleteFile(ctx context.Context, bucket string, objectName string) error

	// EnsureBucketExists ensures that a bucket exists, creates it if it doesn't
	EnsureBucketExists(ctx context.Context, bucketName string) error
}
