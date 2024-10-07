package handlers

import (
	"log"
	"os"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/services"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type UserWithToken struct {
	User  models.User `json:"user"`
	Token string      `json:"token"` // Google OAuth token
}

// GenerateJWTForUser should be imported from the package where it's defined.
func OAuthUserLogin(c *fiber.Ctx) error {
	var request UserWithToken
	jwt_secret := os.Getenv("JWT_SECRET_KEY")

	// Parse the incoming JSON request into the user struct
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Validate the Google OAuth token
	_, err := services.VerifyGoogleOAuthToken(request.Token)
	if err != nil {
		log.Println("Invalid Google OAuth token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid OAuth token",
		})
	}

	// Check if the user exists in the database
	existingUser, isUserExists := services.IsUserExist(request.User.ID)

	if isUserExists {
		// Compare and update user data if necessary
		updatedUser, err := services.UpdateUserNameAndAvatar(existingUser, &request.User)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update user",
			})
		}

		// Generate JWT for the existing user
		token, err := services.GenerateJWTForUser(*updatedUser, jwt_secret)
		if err != nil {
			log.Println("Failed to generate token:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate JWT token",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user":  updatedUser,
			"token": token,
		})
	}

	// Create a new user if it doesn't exist
	createdUser, err := storage.CreateUser(request.User)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Generate JWT for the new user
	token, err := services.GenerateJWTForUser(*createdUser, jwt_secret)
	if err != nil {
		log.Println("Failed to generate token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate JWT token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":  createdUser,
		"token": token,
	})
}
