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
	ID                  uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID              string          `gorm:"not null" json:"user_id"` // Foreign key to User
	User                User            `gorm:"foreignKey:UserID"`       // Belongs to User
	Actors              ActorArray      `gorm:"type:jsonb" json:"actors"`
	NotificationContent string          `json:"notification_content"`
	ReferenceContent    string          `json:"reference_content"`
	ActionType          ActionTypeArray `gorm:"type:jsonb" json:"action_type"`
	ReferenceID         uuid.UUID       `gorm:"type:uuid" json:"reference_id"`
	IsRead              bool            `gorm:"default:false" json:"is_read"`
	CreatedAt           time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
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
