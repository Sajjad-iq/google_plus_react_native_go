package handlers

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreatePost(c *fiber.Ctx) error {
	post := new(models.Post)

	// Parse the incoming JSON request into the post struct
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Generate a new UUID for the post
	post.ID = uuid.New()

	// Create the post in the database
	if err := storage.CreatePost(*post); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create post",
		})
	}

	// Return the created post with a 201 status code
	return c.Status(fiber.StatusCreated).JSON(post)
}
