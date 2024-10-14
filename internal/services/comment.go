package services

import (
	"fmt"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models/requestModels"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/utils"
	"github.com/google/uuid"
)

// DeleteCommentService handles the logic of deleting a comment
func DeleteCommentService(commentID uuid.UUID, userID string) error {
	// Fetch the comment by its ID
	comment, err := storage.FindCommentByID(commentID)
	if err != nil {
		return fmt.Errorf("failed to find comment: %w", err)
	}

	// Ensure the user is the author of the comment or has admin privileges
	if comment.UserID != userID {
		return fmt.Errorf("unauthorized action: user is not the author of the comment")
	}

	// Call the storage function to delete the comment
	if err := storage.DeleteComment(commentID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	// Update the comments counter on the post
	if err := DecrementPostCommentsCounter(comment.PostID); err != nil {
		return fmt.Errorf("failed to update comments counter: %w", err)
	}

	return nil
}

// updateCommentsCounterOnDelete decrements the comments count for a given post ID
func DecrementPostCommentsCounter(postID uuid.UUID) error {
	post, err := storage.GetPostByID(postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	if post.CommentsCount > 0 {
		// Decrement the comments count
		post.CommentsCount--
	}

	// Save the updated post
	if err := database.DB.Save(&post).Error; err != nil {
		return fmt.Errorf("failed to save updated post: %w", err)
	}

	return nil
}

// CreateCommentService creates a new comment on a post
func CreateCommentService(postID uuid.UUID, userID string, commentRequestBody requestModels.CreateCommentRequestBody, lang string) (*models.Comment, error) {
	// Validate comment content
	if err := utils.ValidateCommentContent(commentRequestBody.Content); err != nil {
		return nil, err
	}

	// Fetch user by ID
	commentedUser, err := storage.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Fetch user by ID
	post, err := storage.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	// Create a new comment
	newComment, err := utils.CreateNewComment(post.ID, commentRequestBody, commentedUser)
	if err != nil {
		return nil, err
	}

	// Save the comment
	if err := storage.SaveComment(newComment); err != nil {
		return nil, fmt.Errorf("failed to save comment: %w", err)
	}

	// Increment post comments counter
	if _, err := IncrementPostCommentsCounter(post.ID); err != nil {
		return nil, fmt.Errorf("failed to update comments counter: %w", err)
	}

	// Handle notifications
	if err := handleCommentNotifications(commentRequestBody, commentedUser, *post); err != nil {
		return nil, err
	}

	return newComment, nil
}

// updateCommentsCounter increments the comments count for a given post ID
func IncrementPostCommentsCounter(postID uuid.UUID) (*models.Post, error) {

	post, err := storage.GetPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Increment the comments count
	post.CommentsCount++

	// Save the updated post
	if err := database.DB.Save(&post).Error; err != nil {
		return nil, fmt.Errorf("failed to save updated post: %w", err)
	}

	return post, nil
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
		Limit(limit).
		Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func handleCommentNotifications(commentRequestBody requestModels.CreateCommentRequestBody, commentedUser *models.User, post models.Post) error {
	var actionTypes []string

	// Handle mention notifications
	if len(commentRequestBody.MentionedUsers) > 0 && commentRequestBody.MentionedUsers[0].UserID != "" {
		actionTypes = append(actionTypes, "mention")
		notifyUser, err := storage.FindUserByID(commentRequestBody.MentionedUsers[0].UserID)
		if err != nil {
			return fmt.Errorf("failed to create or update notification: %w", err)
		}

		_, err = CreateOrUpdateNotification(notifyUser, commentedUser.ID, actionTypes, post.ID, commentRequestBody.Content)
		if err != nil {
			return fmt.Errorf("failed to create or update notification: %w", err)
		}
	} else if post.AuthorID != commentedUser.ID {
		// Handle comment notifications to post author
		notifyUser, err := storage.FindUserByID(post.AuthorID)
		if err != nil {
			return fmt.Errorf("failed to create or update notification: %w", err)
		}

		actionTypes = append(actionTypes, "comment")
		_, err = CreateOrUpdateNotification(notifyUser, commentedUser.ID, actionTypes, post.ID, commentRequestBody.Content)
		if err != nil {
			return fmt.Errorf("failed to create or update notification: %w", err)
		}
	}

	return nil
}
