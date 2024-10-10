package storage

import (
	"fmt"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/google/uuid"
)

// FindLikeByUserAndPost fetches a like by user and post from the database
func FindLikeByUserAndPost(postID uuid.UUID, userID string) (*models.Like, error) {
	var existingLike models.Like
	err := database.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&existingLike).Error
	if err != nil {
		return nil, err
	}
	return &existingLike, nil
}

// CreateLike creates a new like in the database
func CreateLike(like *models.Like) error {
	if err := database.DB.Create(like).Error; err != nil {
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

// DeleteLike removes an existing like from the database
func DeleteLike(like *models.Like) error {
	if err := database.DB.Delete(like).Error; err != nil {
		return fmt.Errorf("failed to remove like: %w", err)
	}
	return nil
}

// DeleteLikesByPostID deletes all likes associated with a specific post
func DeleteLikesByPostID(postID uuid.UUID) error {
	if err := database.DB.Where("post_id = ?", postID).Delete(&models.Like{}).Error; err != nil {
		return fmt.Errorf("could not delete likes for post %v: %w", postID, err)
	}
	return nil
}
