package storage

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
)

// SaveComment saves a comment in the database
func SaveComment(comment *models.Comment) error {
	if err := database.DB.Create(comment).Error; err != nil {
		return err
	}
	return nil
}
