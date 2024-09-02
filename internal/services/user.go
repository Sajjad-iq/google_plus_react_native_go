package services

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"gorm.io/gorm"
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

// CompareUserData compares the fields of two user objects and returns a map of the changes.
// It also returns a boolean indicating whether there are any differences.
func CompareUserData(existingUser, requestUser *models.User) (map[string]interface{}, bool) {
	changes := make(map[string]interface{})

	// Compare fields and add to the changes map if they differ
	if requestUser.Username != existingUser.Username {
		changes["Username"] = requestUser.Username
	}
	if requestUser.ProfileAvatar != existingUser.ProfileAvatar {
		changes["ProfileAvatar"] = requestUser.ProfileAvatar
	}

	return changes, len(changes) > 0
}

// UpdateUserChanges applies the changes from requestUser to existingUser if there are differences.
// It returns an error if the update fails.
func UpdateUserChanges(existingUser, requestUser *models.User) error {
	// Compare existing and request user data
	changes, hasChanges := CompareUserData(existingUser, requestUser)
	if !hasChanges {
		// No changes to update
		return nil
	}

	// Apply changes to the existing user
	for key, value := range changes {
		switch key {
		case "Username":
			existingUser.Username = value.(string)
		case "ProfileAvatar":
			existingUser.ProfileAvatar = value.(string)
		}
	}

	// Save updated user to the database
	if err := storage.UpdateUser(*existingUser); err != nil {
		// Return the error if the update fails
		return err
	}

	return nil
}

func UpdateUserAndPosts(db *gorm.DB, user *models.User) error {
	// Begin a transaction
	tx := db.Begin()

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
