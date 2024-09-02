package services

import (
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

func CompereUserData(existingUser, requestUser *models.User) (map[string]interface{}, bool) {
	// Initialize a map to hold the changed fields
	changes := make(map[string]interface{})

	// Compare fields and add to the changes map if they differ
	if requestUser.Username != existingUser.Username {
		changes["Username"] = requestUser.Username
	}

	if requestUser.ProfileAvatar != existingUser.ProfileAvatar {
		changes["ProfileAvatar"] = requestUser.ProfileAvatar
	}

	// Return the changes map and a boolean indicating if there are any changes
	return changes, len(changes) > 0
}

// UpdateUserChanges updates the fields of an existing user with data from a request
// if there are changes. It returns an error if the update fails.
func UpdateUserChanges(existingUser, requestUser *models.User) error {
	// Identify changes between existing user and request user
	_, hasChanges := CompereUserData(existingUser, requestUser)
	if !hasChanges {
		// No changes to update
		return nil
	}

	existingUser.Username = requestUser.Username
	existingUser.ProfileAvatar = requestUser.ProfileAvatar

	// Update the user in the database
	if err := storage.UpdateUser(*existingUser); err != nil {
		// Return the error if the update fails
		return err
	}

	return nil
}
