package storage

import (
	"fmt"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
)

// CreatePost creates a new post in the database
func CreatePost(post models.Post) error {
	// Add database logic here (e.g., GORM or raw SQL)
	// Example:
	if err := database.DB.Create(&post).Error; err != nil {
		return fmt.Errorf("could not create post: %w", err)
	}
	return nil
}
