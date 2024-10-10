package services

import (
	"fmt"
	"log" // Improved error logging
	"time"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/google/uuid"
)

// CreateOrUpdateNotification handles updating or creating a notification
func CreateOrUpdateNotification(userID, actorID string, actionTypes []string, referenceID uuid.UUID, notificationContent string) (*models.Notification, error) {
	// Check for an existing notification
	existingNotification, err := storage.FindNotificationByUserActionAndReference(userID, referenceID)
	if err != nil {
		log.Println("Error checking for existing notification:", err)
		return nil, fmt.Errorf("failed to check for existing notification: %w", err)
	}

	// Fetch or create the actor
	actor, err := FindActor(actorID)
	if err != nil {
		log.Println("Error finding actor:", err)
		return nil, err
	}

	if existingNotification != nil {
		// Update the existing notification
		if err := UpdateExistingNotification(existingNotification, actor, actionTypes); err != nil {
			log.Println("Error updating existing notification:", err)
			return nil, err
		}
		return existingNotification, nil
	}

	// Create a new notification
	return CreateNewNotification(userID, actor, actionTypes, referenceID, notificationContent)
}

// FindActor retrieves actor information and creates an Actor model
func FindActor(actorID string) (models.Actor, error) {
	actorUser, err := storage.FindUserByID(actorID) // Assuming FindUserByID takes a string
	if err != nil {
		log.Println("Error finding actor user:", err)
		return models.Actor{}, fmt.Errorf("failed to find actor user: %w", err)
	}

	actor := models.Actor{
		ID:     actorID,
		Name:   actorUser.Username,
		Avatar: actorUser.ProfileAvatar,
	}
	return actor, nil
}

// UpdateExistingNotification updates an existing notification with the new actor
// UpdateExistingNotification appends new actor and action types if not already present
func UpdateExistingNotification(notification *models.Notification, actor models.Actor, newActionTypes []string) error {
	// Check if the actor is already part of the notification
	actorExists := false
	for _, existingActor := range notification.Actors {
		if existingActor.ID == actor.ID {
			actorExists = true
			break
		}
	}
	if !actorExists {
		notification.Actors = append(notification.Actors, actor)
	}

	// Add any new action types that are not already present
	for _, newActionType := range newActionTypes {
		actionExists := false
		for _, existingActionType := range notification.ActionType {
			if existingActionType == newActionType {
				actionExists = true
				break
			}
		}
		if !actionExists {
			notification.ActionType = append(notification.ActionType, newActionType)
		}
	}

	// Update the timestamp
	notification.UpdatedAt = time.Now()

	// Save the updated notification
	if err := storage.SaveNotification(notification); err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

// CreateNewNotification creates a new notification entry
func CreateNewNotification(userID string, actor models.Actor, actionTypes []string, referenceID uuid.UUID, notificationContent string) (*models.Notification, error) {
	newNotification := models.Notification{
		ID:                  uuid.New(),
		UserID:              userID,
		Actors:              []models.Actor{actor}, // Adding the actor to the Actors array
		ActionType:          actionTypes,           // Using the array of action types
		ReferenceID:         referenceID,
		NotificationContent: notificationContent,
		IsRead:              false,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Save the new notification
	if err := storage.SaveNotification(&newNotification); err != nil {
		log.Println("Error creating new notification:", err)
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &newNotification, nil
}

// DeleteNotificationService deletes a notification by its ID
func DeleteNotificationService(notificationID uuid.UUID) error {
	if err := storage.DeleteNotification(notificationID); err != nil {
		log.Println("Error deleting notification:", err)
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// FetchUserNotificationsService fetches notifications for a user
func FetchUserNotificationsService(userID string, limit int) ([]models.Notification, error) {
	notifications, err := storage.FetchNotificationsByUserID(userID, limit)
	if err != nil {
		log.Println("Error fetching user notifications:", err)
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}
	return notifications, nil
}
