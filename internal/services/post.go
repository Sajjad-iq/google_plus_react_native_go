package services

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
)

func CreatePost(post *models.Post) {
	storage.SavePost(*post)
}

func GetAllPosts() []models.Post {
	return storage.GetAllPosts()
}
