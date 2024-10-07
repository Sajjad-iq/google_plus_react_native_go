package requestModels

type MentionedUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CreateCommentRequestBody struct {
	Content        string          `json:"content"`
	MentionedUsers []MentionedUser `json:"mentioned_users"`
}
