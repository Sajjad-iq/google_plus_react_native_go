package storage

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
)

// CreateUser saves a user to the database and returns an error if it fails.
func CreateUser(user models.User) error {
	if err := database.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// FindUserByID finds a user by their ID and returns the user object and an error if any.
func FindUserByID(id string) (*models.User, error) {
	var user models.User

	// Query the database for the user by ID
	err := database.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		// Return nil user and nil error if no record is found
		return nil, err
	}
	// Return the user and nil error if found successfully
	return &user, nil
}

// UpdateUser updates an existing user in the database and returns an error if it fails.
func UpdateUser(user models.User) error {
	// Use the Save method to update the user record.
	// This method will perform an update if the primary key exists.
	if err := database.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}
