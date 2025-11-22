package entity

import (
	"time"

	"github.com/google/uuid"
)

// ImageStatus represents the status of an image processing
type ImageStatus string

const (
	ImageStatusPending    ImageStatus = "pending"
	ImageStatusProcessing ImageStatus = "processing"
	ImageStatusCompleted  ImageStatus = "completed"
	ImageStatusFailed     ImageStatus = "failed"
)

// Image represents an image entity in the domain
type Image struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	OriginalURL  string      `json:"original_url" db:"original_url"`
	ProcessedURL *string     `json:"processed_url" db:"processed_url"`
	Status       ImageStatus `json:"status" db:"status"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

// IsValid checks if the image status is valid
func (s ImageStatus) IsValid() bool {
	switch s {
	case ImageStatusPending, ImageStatusProcessing, ImageStatusCompleted, ImageStatusFailed:
		return true
	default:
		return false
	}
}

// String returns the string representation of ImageStatus
func (s ImageStatus) String() string {
	return string(s)
}
