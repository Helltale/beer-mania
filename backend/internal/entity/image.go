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
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()" db:"id"`
	OriginalURL  string      `json:"original_url" gorm:"type:varchar(512);not null" db:"original_url"`
	ProcessedURL *string     `json:"processed_url" gorm:"type:varchar(512)" db:"processed_url"`
	Status       ImageStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending','processing','completed','failed')" db:"status"`
	CreatedAt    time.Time   `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP" db:"updated_at"`
}

// TableName specifies the table name for Image
func (Image) TableName() string {
	return "images"
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
