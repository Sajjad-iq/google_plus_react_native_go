package handlers

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/services"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Helper function for checking user authentication
func ValidateRequest(c *fiber.Ctx) (string, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	// Ensure the 'id' exists in the token claims
	id, idOk := claims["id"].(string)
	if !idOk || id == "" {
		return "", fmt.Errorf("unauthorized user")
	}
	return id, nil
}

func DeletePost(c *fiber.Ctx) error {
	// Ensure the user is authenticated
	userID, err := ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	// Get the post ID from the URL
	postID := c.Params("id")
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	// Fetch the post from the database
	post, err := storage.GetPostByID(postUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	// Check if the authenticated user is the author of the post
	if post.AuthorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to delete this post",
		})
	}

	// Delete all associated likes for the post
	if err := storage.DeleteLikesByPostID(postUUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete post likes",
		})
	}

	// Delete all associated comments for the post
	if err := storage.DeleteCommentsByPostID(postUUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete post comments",
		})
	}

	// Delete the post itself
	if err := storage.DeletePost(postUUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete post",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Post deleted successfully by user %s", userID),
	})
}

func GetPostByID(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	// Check if the user has liked the post
	var like models.Like
	err = database.DB.Where("post_id = ? AND user_id = ?", post.ID, userID).First(&like).Error
	if err == nil {
		// User has liked the post
		post.YourLike = true
	} else {
		post.YourLike = false
	}

	// Return the post with the 'YourLike' field included
	return c.Status(fiber.StatusOK).JSON(post)
}

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
	// Toggle the like state based on the user
	liked, err := storage.ToggleLike(post, userID) // Modify to return like status
	if err != nil {
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

func GetPosts(c *fiber.Ctx) error {
	// Ensure the user is authenticated
	userID, err := ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	// Get the limit from query parameters, default to 10 if not provided
	limitParam := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
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

	// Fetch the likes for the current user in bulk
	var likes []models.Like
	err = database.DB.Where("user_id = ?", userID).Find(&likes).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve likes",
		})
	}

	likedPostIDs := make(map[uuid.UUID]bool)
	for _, like := range likes {
		likedPostIDs[like.PostID] = true
	}

	// Set 'YourLike' for each post
	for i := range posts {
		if likedPostIDs[posts[i].ID] {
			posts[i].YourLike = true
		} else {
			posts[i].YourLike = false
		}
	}

	stop := len(posts) < limit

	// Return the posts as JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"stop":  stop,  // true if no more data to load, false otherwise
		"posts": posts, // the fetched posts
	})
}

func CreatePost(c *fiber.Ctx) error {
	// Ensure the user is authenticated
	_, err := ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

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
