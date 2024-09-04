package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreatePost(c *fiber.Ctx) error {
	// Parse form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse form data",
		})
	}

	// Create a new Post instance
	post := new(models.Post)

	// Validate and assign form values to the post
	if authorIDs, ok := form.Value["author_id"]; ok && len(authorIDs) > 0 {
		post.AuthorID = authorIDs[0]
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "author_id is required",
		})
	}

	if authorNames, ok := form.Value["author_name"]; ok && len(authorNames) > 0 {
		post.AuthorName = authorNames[0]
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "author_name is required",
		})
	}

	if bodies, ok := form.Value["body"]; ok && len(bodies) > 0 {
		post.Body = bodies[0]
	}

	if author_avatar, ok := form.Value["author_avatar"]; ok && len(author_avatar) > 0 {
		post.AuthorAvatar = author_avatar[0]
	}

	if shareStates, ok := form.Value["share_state"]; ok && len(shareStates) > 0 {
		post.ShareState = shareStates[0]
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "share_state is required",
		})
	}

	// Generate a new UUID for the post
	post.ID = uuid.New()

	// Handle image upload if present
	if files := form.File["image_url"]; len(files) > 0 {
		// Take the first image if multiple are uploaded
		file := files[0]

		imageUUID := uuid.New().String()
		imagePath := filepath.Join("uploads", imageUUID+filepath.Ext(file.Filename))

		// Ensure the directory structure exists
		if err := os.MkdirAll(filepath.Dir(imagePath), os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create directory",
			})
		}

		// Save the image
		if err := c.SaveFile(file, imagePath); err != nil {
			fmt.Printf("Error saving file: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save image",
			})
		}

		// Set the ImageURL field to the saved image path
		post.ImageURL = imagePath
	} else {
		// If no image is provided, set ImageURL to an empty string
		post.ImageURL = ""
	}

	// Create the post in the database
	if err := storage.CreatePost(*post); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create post",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}
