package storage

import (
	"fmt"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
)

// CreatePost creates a new post in the database

func ToggleLike(post *models.Post, userID string) error {
	// Check if user has already liked the post
	var existingLike models.Like
	err := database.DB.Where("post_id = ? AND user_id = ?", post.ID, userID).First(&existingLike).Error

	if err == nil {
		// User has already liked the post, remove like
		if err := database.DB.Delete(&existingLike).Error; err != nil {
			return fmt.Errorf("failed to remove like: %w", err)
		}
		post.LikesCount--
	} else {
		// User has not liked the post, add a new like
		newLike := models.Like{
			UserID: userID,
			PostID: post.ID,
		}
		if err := database.DB.Create(&newLike).Error; err != nil {
			return fmt.Errorf("failed to add like: %w", err)
		}
		post.LikesCount++
	}

	// Update post in the database
	return UpdatePost(post)
}
