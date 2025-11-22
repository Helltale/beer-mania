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
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()" db:"id"`
	ImageID      uuid.UUID  `json:"image_id" gorm:"type:uuid;not null;index" db:"image_id"`
	Status       TaskStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending';index;check:status IN ('pending','processing','completed','failed')" db:"status"`
	ErrorMessage *string    `json:"error_message" gorm:"type:text" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;index" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP" db:"updated_at"`
}

// TableName specifies the table name for ProcessingTask
func (ProcessingTask) TableName() string {
	return "processing_tasks"
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
