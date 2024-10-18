package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// MentionedUser represents a user who is mentioned in the comment
type MentionedUser struct {
	UserID   string `gorm:"not null" json:"user_id"` // User ID of the mentioned user
	UserName string `json:"user_name"`               // Username of the mentioned user
}

// MentionedUserArray represents an array of MentionedUser objects
type MentionedUserArray []MentionedUser

// Comment represents a comment on a post
type Comment struct {
	ID             uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	PostID         uuid.UUID          `gorm:"type:uuid;not null" json:"post_id"` // Foreign key to Post
	Post           Post               `gorm:"foreignKey:PostID"`                 // Belongs to Post
	UserID         string             `gorm:"not null" json:"user_id"`           // Foreign key to User
	User           User               `gorm:"foreignKey:UserID"`                 // Belongs to User
	AuthorName     string             `json:"author_name"`
	AuthorAvatar   string             `json:"author_avatar"`
	Content        string             `gorm:"type:text;not null" json:"content"`
	CreatedAt      time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	MentionedUsers MentionedUserArray `gorm:"type:jsonb" json:"mentioned_users"` // Array of mentioned users stored as JSONB
}

// Scan implements the sql.Scanner interface for MentionedUserArray
func (m *MentionedUserArray) Scan(value interface{}) error {
	if value == nil {
		*m = MentionedUserArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan mentioned users: expected []byte")
	}

	var mentionedUsers MentionedUserArray
	if err := json.Unmarshal(bytes, &mentionedUsers); err != nil {
		return errors.New("failed to unmarshal mentioned users: " + err.Error())
	}

	*m = mentionedUsers
	return nil
}

// Value implements the driver.Valuer interface for MentionedUserArray
func (m MentionedUserArray) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}
