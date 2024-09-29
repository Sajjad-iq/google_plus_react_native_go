package storage

import (
	"fmt"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/google/uuid"
)

// SaveComment saves a comment in the database
func SaveComment(comment *models.Comment) error {
	if err := database.DB.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCommentsByPostID(postID uuid.UUID) error {
	if err := database.DB.Where("post_id = ?", postID).Delete(&models.Comment{}).Error; err != nil {
		return fmt.Errorf("could not delete comments for post %v: %w", postID, err)
	}
	return nil
}

// FindCommentByID retrieves a comment by its ID
func FindCommentByID(commentID uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	if err := database.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteComment removes a comment from the database by its ID
func DeleteComment(commentID uuid.UUID) error {
	if err := database.DB.Where("id = ?", commentID).Delete(&models.Comment{}).Error; err != nil {
		return err
	}
	return nil
}
