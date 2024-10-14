package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/services"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// FetchNotificationsHandler handles the request for fetching user notifications
func FetchNotificationsHandler(c *fiber.Ctx) error {
	// Extract the userID from the request context (assuming it's set by middleware)
	userID, err := ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	// Access the Accept-Language header
	lang := c.Get("Accept-Language", "en") // Default to "en" if not set

	// Optional: parse the 'limit' query parameter (default to 10 if not provided)
	limitQuery := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit parameter",
		})
	}

	// Fetch the notifications using the service function
	notifications, err := services.FetchUserNotificationsService(userID, limit, lang) // Pass lang to the service
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notifications",
		})
	}

	// Respond with the notifications
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"notifications": notifications,
	})
}

func MarkNotificationsAsReadHandler(c *fiber.Ctx) error {
	// Extract the userID from the request context (assuming it's set by middleware)
	_, err := ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	notificationID := c.Params("id")
	uuid, err := uuid.Parse(notificationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	// Fetch the notifications using the service function
	err = storage.MarkNotificationAsRead(uuid) // Pass lang to the service
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notifications",
		})
	}

	// Respond with the notifications
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Notification marked as read",
	})
}
