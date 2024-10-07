package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Comment struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"` // Unique identifier for the comment
	PostID         uuid.UUID      `gorm:"type:uuid;not null" json:"post_id"`              // Foreign key, references the post being commented on
	UserID         string         `gorm:"not null" json:"user_id"`                        // Foreign key, references the user making the comment
	AuthorName     string         ` json:"author_name"`                                   // Foreign key, references the user making the comment
	AuthorAvatar   string         ` json:"author_avatar"`                                 // Foreign key, references the user making the comment
	Content        string         `gorm:"type:text;not null" json:"content"`              // Content of the comment
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`               // Timestamp when the comment was created
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`               // Timestamp when the comment was last updated
	MentionedUsers pq.StringArray `gorm:"type:text[]" json:"mentioned_users"`
}
