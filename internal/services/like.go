package services

import (
	"fmt"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
)

// ToggleLike handles the logic for liking or unliking a post and manages notifications
func ToggleLike(post *models.Post, userID string, lang string) (bool, error) {
	// Check if the user has already liked the post
	existingLike, err := storage.FindLikeByUserAndPost(post.ID, userID)
	if err == nil && existingLike != nil {
		// User has already liked the post, remove like
		if err := storage.DeleteLike(existingLike); err != nil {
			return false, fmt.Errorf("failed to remove like: %w", err)
		}

		// Decrement the like count
		post.LikesCount--

		// Update the post in the database
		if err := storage.UpdatePost(post); err != nil {
			return false, fmt.Errorf("failed to update post after removing like: %w", err)
		}

		return false, nil // Returning false to indicate the post is now unliked
	}

	// If the user hasn't liked the post yet, add a new like
	newLike := models.Like{
		UserID: userID,
		PostID: post.ID,
	}
	if err := storage.CreateLike(&newLike); err != nil {
		return false, fmt.Errorf("failed to add like: %w", err)
	}

	// Increment the like count
	post.LikesCount++

	// Update the post in the database
	if err := storage.UpdatePost(post); err != nil {
		return false, fmt.Errorf("failed to update post after adding like: %w", err)
	}

	if post.AuthorID != userID {
		// Create or update a notification for the post like
		actionTypes := []string{"like"} // Define the action type as an array of strings
		_, err = CreateOrUpdateNotification(post.AuthorID, userID, actionTypes, post.ID, post.Body, lang)
		if err != nil {
			return false, fmt.Errorf("failed to create or update notification: %w", err)
		}
	}

	return true, nil // Returning true to indicate the post is now liked
}
