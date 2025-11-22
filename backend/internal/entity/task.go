package entity

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the status of a processing task
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

// ProcessingTask represents a processing task entity in the domain
type ProcessingTask struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ImageID      uuid.UUID  `json:"image_id" db:"image_id"`
	Status       TaskStatus `json:"status" db:"status"`
	ErrorMessage *string    `json:"error_message" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// IsValid checks if the task status is valid
func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusPending, TaskStatusProcessing, TaskStatusCompleted, TaskStatusFailed:
		return true
	default:
		return false
	}
}

// String returns the string representation of TaskStatus
func (s TaskStatus) String() string {
	return string(s)
}
