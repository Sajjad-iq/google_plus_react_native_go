package models

type Post struct {
	Id            string `json:"id"`
	Author        string `json:"author"`
	ShareState    string `json:"shareState"`
	Likes         int    `json:"likes"`
	CommentsCount int    `json:"commentsCount"`
	Body          string `json:"body"`
	Image         string `json:"imageContent"`
	Date          string `json:"date"`
}
