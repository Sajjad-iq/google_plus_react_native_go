package utils

import (
	"fmt"
	"time"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models/requestModels"
	"github.com/google/uuid"
)

func ValidateCommentContent(content string) error {
	if content == "" {
		return fmt.Errorf("comment content cannot be empty")
	}
	return nil
}

func CreateNewComment(postID uuid.UUID, commentRequestBody requestModels.CreateCommentRequestBody, user *models.User) (*models.Comment, error) {
	var mentionedUser models.MentionedUser
	if len(commentRequestBody.MentionedUsers) > 0 && commentRequestBody.MentionedUsers[0].UserID != "" {
		mentionedUser = commentRequestBody.MentionedUsers[0]
	}

	newComment := &models.Comment{
		ID:             uuid.New(),
		PostID:         postID,
		UserID:         user.ID,
		Content:        commentRequestBody.Content,
		MentionedUsers: []models.MentionedUser{mentionedUser},
		AuthorName:     user.Username,
		AuthorAvatar:   user.ProfileAvatar,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return newComment, nil
}
