package storage

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
)

func CreatePost(post models.Post) error {
	if err := database.DB.Create(&post).Error; err != nil {
		return err
	}
	return nil
}
