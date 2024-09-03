package services

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
)

// IsUserExist checks if a user with the given ID exists in the database.
// It returns the user data and true if the user exists, otherwise nil and false.

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
func UpdateUserNameAndAvatar(existingUser, requestUser *models.User) (*models.User, error) {
	// Compare existing and request user data
	changes, hasChanges := CompareUserData(existingUser, requestUser)
	if !hasChanges {
		// No changes to update, return the existing user
		return existingUser, nil
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

	// Save updated user to the database and update the user's posts
	if err := UpdateUserNameAndAvatarForEachBelongingTable(existingUser); err != nil {
		// Return the error if the update fails
		return nil, err
	}

	// Return the updated user
	return existingUser, nil
}
