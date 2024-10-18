package models

import (
	"time"
)

type User struct {
	ID            string    `json:"id" gorm:"primaryKey;type:numeric"`
	Username      string    `json:"username" gorm:"not null"`
	Email         string    `json:"email" gorm:"unique;not null"`
	ProfileAvatar string    `json:"profile_avatar"`
	ProfileCover  string    `json:"profile_cover"`
	Bio           string    `json:"bio"`
	PushToken     string    `json:"push_token"`
	UserLang      string    `json:"user_lang" gorm:"default:'en'"`
	Status        string    `json:"status" gorm:"default:'active'"`
	Role          string    `json:"role" gorm:"default:'user'"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Posts         []Post         `gorm:"foreignKey:AuthorID" json:"posts"`       // One-to-many (User -> Posts)
	Comments      []Comment      `gorm:"foreignKey:UserID" json:"comments"`      // One-to-many (User -> Comments)
	Likes         []Like         `gorm:"foreignKey:UserID" json:"likes"`         // One-to-many (User -> Likes)
	Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications"` // One-to-many (User -> Notifications)
}
