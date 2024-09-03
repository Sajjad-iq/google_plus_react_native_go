package services

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
)

// IsUserExist checks if a user with the given ID exists in the database.
// It returns the user data and true if the user exists, otherwise nil and false.
func IsUserExist(userID string) (*models.User, bool) {
	user, err := storage.FindUserByID(userID)
	if err != nil {
		// If an error occurs (e.g., user not found), return nil and false
		return nil, false
	}
	// If no error occurs, return the user and true
	return user, true
}

// this function used only in case the user change the name or avatar
func UpdateUserNameAndAvatarForEachBelongingTable(user *models.User) error {
	// Begin a transaction
	tx := database.DB.Begin()

	// Rollback the transaction in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update the user
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update all posts belonging to this user
	if err := tx.Model(&models.Post{}).Where("author_id = ?", user.ID).
		Updates(map[string]interface{}{
			"author_name":   user.Username,
			"author_avatar": user.ProfileAvatar,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}
