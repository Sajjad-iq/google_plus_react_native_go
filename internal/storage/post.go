package storage

import "github.com/Sajjad-iq/google_plus_react_native_go/internal/models"

var posts []models.Post

func SavePost(post models.Post) {
	posts = append(posts, post)
}

func GetAllPosts() []models.Post {
	return posts
}
