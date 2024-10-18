package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"post_id"` // Foreign key to Post
	Post      Post      `gorm:"foreignKey:PostID"`                 // Belongs to Post
	UserID    string    `gorm:"not null" json:"user_id"`           // Foreign key to User
	User      User      `gorm:"foreignKey:UserID"`                 // Belongs to User
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`  // Timestamp when the like was created
}
