package storage

import (
	"fmt"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/google/uuid"
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

func GetPosts(limit int) ([]models.Post, error) {
	var posts []models.Post

	// Fetch posts from the database, ordered by 'created_at' field in descending order
	if err := database.DB.Order("created_at DESC").Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostByID(id uuid.UUID) (*models.Post, error) {
	var post models.Post
	if err := database.DB.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func UpdatePost(post *models.Post) error {
	return database.DB.Save(post).Error
}
func DeletePost(id uuid.UUID) error {
	return database.DB.Delete(&models.Post{}, "id = ?", id).Error
}
