package services

import (
	"fmt"
	"time"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/google/uuid"
)

// CreateCommentService handles the logic of creating a new comment
func CreateCommentService(postID uuid.UUID, userID string, content string) (*models.Comment, error) {
	// Validate content is not empty
	if content == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}

	user, err := storage.FindUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Create the new comment
	newComment := models.Comment{
		ID:           uuid.New(),
		PostID:       postID,
		UserID:       userID,
		Content:      content,
		AuthorName:   user.Username,
		AuthorAvatar: user.ProfileAvatar,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Call the storage function to save the comment
	if err := storage.SaveComment(&newComment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Update the comments counter on the post
	if err := updateCommentsCounter(postID); err != nil {
		return nil, fmt.Errorf("failed to update comments counter: %w", err)
	}

	return &newComment, nil
}

// updateCommentsCounter increments the comments count for a given post ID
func updateCommentsCounter(postID uuid.UUID) error {

	post, err := storage.GetPostByID(postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Increment the comments count
	post.CommentsCount++

	// Save the updated post
	if err := database.DB.Save(&post).Error; err != nil {
		return fmt.Errorf("failed to save updated post: %w", err)
	}

	return nil
}

// FetchCommentsService retrieves comments for a given post ID with a limit
func FetchCommentsService(postID uuid.UUID, limit int) ([]models.Comment, error) {
	var comments []models.Comment

	// Validate that limit is greater than zero
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than zero")
	}

	// Fetch comments from the database using the postID, ordered by latest first, with a limit
	if err := database.DB.Where("post_id = ?", postID).
		Order("created_at DESC"). // Order by created_at in descending order
		Limit(limit).
		Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}
