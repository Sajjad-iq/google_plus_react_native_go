package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Actor struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type ActorArray []Actor
type ActionTypeArray []string

// Notification represents a notification in the social media app
type Notification struct {
	ID                  uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"` // Unique identifier for the notification
	UserID              string          `gorm:"not null" json:"user_id"`                        // User receiving the notification
	Actors              ActorArray      `gorm:"type:jsonb" json:"actors"`
	NotificationContent string          `json:"notification_content"`             // Notification content
	ReferenceContent    string          `json:"reference_content"`                // Notification content
	ActionType          ActionTypeArray `gorm:"type:jsonb" json:"action_type"`    // Array of action types (e.g., 'like', 'comment', 'follow', 'mention')
	ReferenceID         uuid.UUID       `gorm:"type:uuid" json:"reference_id"`    // ID referencing the related entity (e.g., post, comment, or user)
	IsRead              bool            `gorm:"default:false" json:"is_read"`     // Whether the notification has been read
	CreatedAt           time.Time       `gorm:"autoCreateTime" json:"created_at"` // Timestamp when the notification was created
	UpdatedAt           time.Time       `gorm:"autoUpdateTime" json:"updated_at"` // Timestamp when the notification was last updated
}

// Scan implements the sql.Scanner interface for ActorArray
func (a *ActorArray) Scan(value interface{}) error {
	if value == nil {
		*a = ActorArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan actors: expected []byte")
	}

	var actors ActorArray
	if err := json.Unmarshal(bytes, &actors); err != nil {
		return errors.New("failed to unmarshal actors: " + err.Error())
	}

	*a = actors
	return nil
}

// Value implements the driver.Valuer interface for ActorArray
func (a ActorArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface for ActionTypeArray
func (a *ActionTypeArray) Scan(value interface{}) error {
	if value == nil {
		*a = ActionTypeArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan action types: expected []byte")
	}

	var actionTypes ActionTypeArray
	if err := json.Unmarshal(bytes, &actionTypes); err != nil {
		return errors.New("failed to unmarshal action types: " + err.Error())
	}

	*a = actionTypes
	return nil
}

// Value implements the driver.Valuer interface for ActionTypeArray
func (a ActionTypeArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}
