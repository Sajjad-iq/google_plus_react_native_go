package handlers

import (
	"mime/multipart"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/services"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
)

func CreatePost(c *fiber.Ctx) error {
	// Parse form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse form data",
		})
	}

	// Validate form and create post struct
	post, err := services.ValidatePostForm(form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Handle image upload if present
	var fileHeader *multipart.FileHeader
	if files := form.File["image_url"]; len(files) > 0 {
		fileHeader = files[0]
		imagePath, err := services.SaveImage(fileHeader)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save image",
			})
		}

		// Set the image URL in the post struct
		post.ImageURL = imagePath
	} else {
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
