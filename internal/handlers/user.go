package handlers

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/services"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
)

func FilterUserLogin(c *fiber.Ctx) error {
	var user models.User

	// Parse the incoming JSON request into the user struct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Check if the user exists in the database
	existingUser, isUserExists := services.IsUserExist(user.ID)

	if isUserExists {
		// Compare and update user data if necessary
		_, hasChanges := services.CompareUserData(existingUser, &user)
		if hasChanges {
			// Apply the changes to the existing user and update the database
			if err := storage.UpdateUser(*existingUser); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to update user data",
				})
			}
		}
		return c.Status(fiber.StatusOK).JSON(existingUser)
	}

	// Create a new user if it doesn't exist
	if err := storage.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
