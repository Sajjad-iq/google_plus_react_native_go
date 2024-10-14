package storage

import (
	"fmt"
	"log" // Improved error logging

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SaveNotification saves a notification in the database
func SaveNotification(notification *models.Notification) error {
	if err := database.DB.Save(notification).Error; err != nil {
		log.Println("Error saving notification:", err)
		return err
	}
	return nil
}

// MarkNotificationAsRead updates a notification to mark it as read
func MarkNotificationAsRead(notificationID uuid.UUID) error {
	notification := &models.Notification{}
	if err := database.DB.First(notification, "id = ?", notificationID).Error; err != nil {
		log.Println("Error finding notification to mark as read:", err)
		return fmt.Errorf("notification not found: %w", err)
	}
	notification.IsRead = true
	if err := database.DB.Save(notification).Error; err != nil {
		log.Println("Error updating notification as read:", err)
		return fmt.Errorf("failed to update notification: %w", err)
	}
	return nil
}

// DeleteNotification deletes a notification by its ID
func DeleteNotification(notificationID uuid.UUID) error {
	if err := database.DB.Where("id = ?", notificationID).Delete(&models.Notification{}).Error; err != nil {
		log.Println("Error deleting notification:", err)
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// FetchNotificationsByUserID retrieves notifications for a specific user
func FetchNotificationsByUserID(userID string, limit int) ([]models.Notification, error) {
	var notifications []models.Notification

	if err := database.DB.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Find(&notifications).Error; err != nil {
		log.Println("Error fetching notifications:", err)
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}

	return notifications, nil
}

// FindNotificationByUserActionAndReference checks if a notification exists for a user, action, and reference
func FindNotificationByUserActionAndReference(userID string, referenceID uuid.UUID) (*models.Notification, error) {
	var notification models.Notification

	err := database.DB.Where("user_id = ? AND reference_id = ?", userID, referenceID).
		First(&notification).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No notification found
		}
		log.Println("Error finding notification by user action and reference:", err)
		return nil, err // Other errors
	}

	return &notification, nil
}
