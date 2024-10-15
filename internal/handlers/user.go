package handlers

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
)

// UpdatePushTokenHandler handles the request for updating a user's push token.
func UpdatePushTokenHandler(c *fiber.Ctx) error {
	// Extract user ID from the request parameters
	userID, err := ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	// Extract the new push token from the request body
	var requestBody struct {
		PushToken string `json:"push_token"`
		UserLang  string `json:"user_lang"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Find the user by ID
	user, err := storage.FindUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update the user's push token
	user.PushToken = requestBody.PushToken
	user.UserLang = requestBody.UserLang

	// Save the updated user record in the database
	if err := storage.UpdateUser(*user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update push token",
		})
	}

	// Respond with a success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Push token updated successfully",
	})
}

func GetTheUser(c *fiber.Ctx) error {
	// Ensure the user is authenticated
	_, err := ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	// Get the post ID from the URL parameters
	requestedUserID := c.Params("id")

	// Fetch the post from the database
	user, err := storage.FindUserByID(requestedUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}
