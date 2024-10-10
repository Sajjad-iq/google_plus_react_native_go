package handlers

import (
	"fmt"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/services"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func LikePost(c *fiber.Ctx) error {
	// Ensure the user is authenticated
	userID, err := ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	// Get the post ID from the URL parameters
	postID := c.Params("id")
	uuid, err := uuid.Parse(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	// Fetch the post from the database
	post, err := storage.GetPostByID(uuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find post",
		})
	}

	// Toggle the like state based on the user
	liked, err := services.ToggleLike(post, userID) // Modify to return like status
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update like status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Post like state updated successfully",
		"likes_count": post.LikesCount,
		"liked":       liked,
	})
}
