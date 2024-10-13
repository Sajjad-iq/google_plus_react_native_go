package services

import (
	"fmt"
	"time"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models/requestModels"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
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

func CreateCommentService(postID uuid.UUID, userID string, commentRequestBody requestModels.CreateCommentRequestBody, lang string) (*models.Comment, error) {
	// Validate content is not empty
	if commentRequestBody.Content == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}

	user, err := storage.FindUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	var MentionedUser models.MentionedUser

	if commentRequestBody.MentionedUsers[0].UserID != "" {
		MentionedUser = commentRequestBody.MentionedUsers[0]
	}

	// Create the new comment
	newComment := models.Comment{
		ID:             uuid.New(),
		PostID:         postID,
		UserID:         userID,
		Content:        commentRequestBody.Content,
		MentionedUsers: []models.MentionedUser{MentionedUser},
		AuthorName:     user.Username,
		AuthorAvatar:   user.ProfileAvatar,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Call the storage function to save the comment
	if err := storage.SaveComment(&newComment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Update the comments counter on the post
	updatedPost, err := IncrementPostCommentsCounter(postID)

	if err != nil {
		return nil, fmt.Errorf("failed to update comments counter: %w", err)
	}

	var actionTypes []string

	if commentRequestBody.MentionedUsers[0].UserID != "" {
		actionTypes = append(actionTypes, "mention")
		_, err = CreateOrUpdateNotification(commentRequestBody.MentionedUsers[0].UserID, userID, actionTypes, postID, commentRequestBody.Content, lang)
		if err != nil {
			return nil, fmt.Errorf("failed to create or update notification: %w", err)
		}
	} else {
		if userID != updatedPost.AuthorID {
			actionTypes = append(actionTypes, "comment")
			_, err = CreateOrUpdateNotification(updatedPost.AuthorID, userID, actionTypes, postID, commentRequestBody.Content, lang)
			if err != nil {
				return nil, fmt.Errorf("failed to create or update notification: %w", err)
			}
		}
	}

	return &newComment, nil
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
