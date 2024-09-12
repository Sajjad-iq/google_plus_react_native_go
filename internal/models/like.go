package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"post_id"` // References the post that is liked
	UserID    string    `gorm:"not null" json:"user_id"`           // References the user who liked the post
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`  // Timestamp when the like was created
}
