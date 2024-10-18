package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Post struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	AuthorID       string         `json:"author_id"`           // Foreign key to User
	Author         User           `gorm:"foreignKey:AuthorID"` // Belongs to User
	AuthorName     string         `json:"author_name"`
	AuthorAvatar   string         `json:"author_avatar"`
	Body           string         `json:"body"`
	ImageURL       string         `json:"image_url"`
	ShareState     string         `gorm:"default:Public" json:"share_state"`
	LikesCount     int            `gorm:"default:0" json:"likes_count"`
	CommentsCount  int            `gorm:"default:0" json:"comments_count"`
	Hashtags       pq.StringArray `gorm:"type:text[]" json:"hashtags"`
	MentionedUsers pq.Int32Array  `gorm:"type:int[]" json:"mentioned_users"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	YourLike       bool           `json:"your_like"` // Computed at runtime

	// Relationships
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments"` // One-to-many (Post -> Comments)
	Likes    []Like    `gorm:"foreignKey:PostID" json:"likes"`    // One-to-many (Post -> Likes)
}
