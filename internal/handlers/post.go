package handlers

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/services"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func DeletePost(c *fiber.Ctx) error {
	// Get the post ID
	postID := c.Params("id")
	uuid, err := uuid.Parse(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	// Delete the post
	if err := storage.DeletePost(uuid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete post",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}
func GetPostByID(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	// Return the post as JSON
	return c.Status(fiber.StatusOK).JSON(post)
}

func LikePost(c *fiber.Ctx) error {
	// Get the user from the JWT token
	userID := c.Locals("user").(string)

	// Ensure the user is authenticated
	if userID == "" {
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
	if err := storage.ToggleLike(post, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update like status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Post like state updated successfully",
		"likes_count": post.LikesCount,
	})
}
func GetPosts(c *fiber.Ctx) error {
	// Get the limit from query parameters, default to 10 if not provided
	limitParam := c.Query("limit", "10")

	// Convert the limit from string to integer
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		// Handle invalid limit input
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit parameter",
		})
	}

	// Fetch posts from the database with the specified limit
	posts, err := storage.GetPosts(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve posts",
		})
	}

	// Return posts as JSON
	return c.Status(fiber.StatusOK).JSON(posts)
}

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

		// Set the full image URL in the post struct
		post.ImageURL = fmt.Sprintf("%s/uploads/%s", c.BaseURL(), filepath.Base(imagePath))
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
