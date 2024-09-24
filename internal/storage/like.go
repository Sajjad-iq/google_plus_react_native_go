package storage

import (
	"fmt"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
)

// CreatePost creates a new post in the database

func ToggleLike(post *models.Post, userID string) (bool, error) {
	var existingLike models.Like

	// Check if the user has already liked the post
	err := database.DB.Where("post_id = ? AND user_id = ?", post.ID, userID).First(&existingLike).Error
	if err == nil {
		// User has already liked the post, remove like
		if err := database.DB.Delete(&existingLike).Error; err != nil {
			return false, fmt.Errorf("failed to remove like: %w", err)
		}

		// Decrement the like count
		post.LikesCount--

		// Update the post in the database
		if err := UpdatePost(post); err != nil {
			return false, fmt.Errorf("failed to update post after removing like: %w", err)
		}

		return false, nil // Returning false to indicate the post is now unliked
	}

	// If the user hasn't liked the post yet, add a new like
	newLike := models.Like{
		UserID: userID,
		PostID: post.ID,
	}
	if err := database.DB.Create(&newLike).Error; err != nil {
		return false, fmt.Errorf("failed to add like: %w", err)
	}

	// Increment the like count
	post.LikesCount++

	// Update the post in the database
	if err := UpdatePost(post); err != nil {
		return false, fmt.Errorf("failed to update post after adding like: %w", err)
	}

	return true, nil // Returning true to indicate the post is now liked
}
