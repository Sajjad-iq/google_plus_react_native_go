package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models/requestModels"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DeleteComment deletes a comment from a post
func DeleteComment(c *fiber.Ctx) error {
	// Ensure the user is authenticated
	userID, err := ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	// Get the comment ID from the URL
	commentID := c.Params("id")
	commentUUID, err := uuid.Parse(commentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	// Delete the comment
	err = services.DeleteCommentService(commentUUID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete comment",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Comment deleted successfully by user %s", userID),
	})
}

// CreateComment handles creating a new comment for a post
func CreateComment(c *fiber.Ctx) error {
	// Ensure the user is authenticated
	userID, err := ValidateRequest(c) // Assuming you have a method to validate the user
	if err != nil {
		log.Println("Error: Unauthorized user -", err) // Log the error
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	// Get the post ID from the request parameters
	postID := c.Params("id")
	uuidPostID, err := uuid.Parse(postID)
	if err != nil {
		log.Println("Error: Invalid post ID -", err) // Log the error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	// Parse the request body to get the comment content
	var requestBody requestModels.CreateCommentRequestBody

	if err := c.BodyParser(&requestBody); err != nil {
		log.Println("Error: Invalid request body -", err) // Log the error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Call the service to handle the comment creation
	comment, err := services.CreateCommentService(uuidPostID, userID, requestBody)
	if err != nil {
		log.Println("Error: Failed to create comment -", err) // Log the error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// FetchComments handles fetching all comments for a post with an optional limit
func FetchComments(c *fiber.Ctx) error {
	// Ensure the user is authenticated
	_, err := ValidateRequest(c) // Assuming you have a method to validate the user
	if err != nil {
		log.Println("Error: Unauthorized user -", err) // Log the error
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	// Get the post ID from the request parameters
	postID := c.Params("id")
	uuidPostID, err := uuid.Parse(postID)
	if err != nil {
		log.Println("Error: Invalid post ID -", err) // Log the error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	// Get the limit query parameter, if provided
	limitParam := c.Query("limit", "10") // Default to 10 comments if limit not provided
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		log.Println("Error: Invalid limit parameter -", err) // Log the error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit parameter",
		})
	}

	// Call the service to fetch the comments for the post with the limit
	comments, err := services.FetchCommentsService(uuidPostID, limit)
	if err != nil {
		log.Println("Error: Failed to fetch comments for post ID:", postID, "-", err) // Log the error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch comments",
		})
	}

	// Determine if "stop" should be true (i.e., when the number of comments fetched is less than the limit)
	stop := len(comments) < limit

	// Return the list of comments and stop flag in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"stop":     stop,     // true if no more data to load, false otherwise
		"comments": comments, // the fetched comments
	})
}
