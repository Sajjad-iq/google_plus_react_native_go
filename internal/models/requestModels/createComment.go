package requestModels

import "github.com/Sajjad-iq/google_plus_react_native_go/internal/models"

type CreateCommentRequestBody struct {
	Content        string                 `json:"content"`
	MentionedUsers []models.MentionedUser `json:"mentioned_users"`
}
