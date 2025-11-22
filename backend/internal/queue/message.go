package queue

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ProcessingMessage struct {
	TaskID  uuid.UUID `json:"task_id"`
	ImageID uuid.UUID `json:"image_id"`
}

func (m *ProcessingMessage) Validate() error {
	if m.TaskID == uuid.Nil {
		return errors.New("task_id cannot be nil")
	}
	if m.ImageID == uuid.Nil {
		return errors.New("image_id cannot be nil")
	}
	return nil
}

func (m *ProcessingMessage) Marshal() ([]byte, error) {
	if err := m.Validate(); err != nil {
		return nil, fmt.Errorf("message validation failed: %w", err)
	}
	return json.Marshal(m)
}

func UnmarshalProcessingMessage(data []byte) (*ProcessingMessage, error) {
	var msg ProcessingMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal processing message: %w", err)
	}

	if err := msg.Validate(); err != nil {
		return nil, fmt.Errorf("unmarshaled message validation failed: %w", err)
	}

	return &msg, nil
}
