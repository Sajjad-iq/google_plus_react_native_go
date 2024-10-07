package storage

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
)

// CreateUser inserts a new user record into the database.
// It returns an error if the operation fails.
func CreateUser(user models.User) (*models.User, error) {
	// Attempt to create the user in the database
	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err // Return error if the creation fails
	}

	// Return the created user and nil error if successful
	return &user, nil
}

// FindUserByID retrieves a user by their ID from the database.
// It returns a pointer to the user object and an error, if any occurred during the operation.
func FindUserByID(id string) (*models.User, error) {
	var user models.User

	// Query the database for the user by ID
	err := database.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		// If the user is not found or another error occurs, return nil and the error
		return nil, err
	}

	// Return the user and nil error if found successfully
	return &user, nil
}

// UpdateUser updates an existing user record in the database.
// It returns an error if the operation fails.
func UpdateUser(user models.User) error {
	// Save the user record to the database. This will update the existing record if the primary key exists.
	if err := database.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}
