package handlers

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/services"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
)

func FilterUserLogin(c *fiber.Ctx) error {
	user := new(models.User)

	// Parse the incoming JSON request into the user struct
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Check if the user exists
	existsUserData, isUserExists := services.IsUserExist(user.ID)

	if !isUserExists {
		// Create the user if it doesn't exist
		if err := storage.CreateUser(*user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create user",
			})
		}
	} else {
		// Compare and update user data if necessary
		if err := services.UpdateUserChanges(existsUserData, user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not update user",
			})
		}
	}

	if isUserExists {
		return c.Status(fiber.StatusCreated).JSON(existsUserData)
	} else {
		return c.Status(fiber.StatusCreated).JSON(user)
	}
}
