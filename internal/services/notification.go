package services

import (
	"fmt"
	"log" // Improved error logging
	"time"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/utils"
	"github.com/google/uuid"
)

// CreateOrUpdateNotification handles updating or creating a notification
func CreateOrUpdateNotification(notifyUser *models.User, actorID string, actionTypes []string, referenceID uuid.UUID, referenceContent string) (*models.Notification, error) {
	// Check for an existing notification
	existingNotification, err := storage.FindNotificationByUserActionAndReference(notifyUser.ID, referenceID)
	if err != nil {
		log.Println("Error checking for existing notification:", err)
		return nil, fmt.Errorf("failed to check for existing notification: %w", err)
	}

	// Fetch or create the actor
	actor, err := CreateActor(actorID)
	if err != nil {
		log.Println("Error finding actor:", err)
		return nil, err
	}

	var notification *models.Notification

	if existingNotification != nil {
		// Update the existing notification
		if err := updateExistingNotification(existingNotification, actor, actionTypes, referenceContent); err != nil {
			log.Println("Error updating existing notification:", err)
			return nil, err
		}
		notification = existingNotification
	} else {
		// Create a new notification
		var err error
		notification, err = createNewNotification(notifyUser.ID, actor, actionTypes, referenceID, referenceContent)
		if err != nil {
			log.Println("Error creating new notification:", err)
			return nil, err
		}
	}

	// Send notification using Expo Push Notification API
	if err := utils.SendPushNotification(notifyUser, notification); err != nil {
		log.Println("Error sending push notification:", err)
		return nil, err
	}

	return notification, nil
}

// FindActor retrieves actor information and creates an Actor model
func CreateActor(actorID string) (models.Actor, error) {
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
func updateExistingNotification(notification *models.Notification, actor models.Actor, newActionTypes []string, ReferenceContent string) error {
	// Check if the actor is already part of the notification
	actorExists := false
	for i, existingActor := range notification.Actors {
		if existingActor.ID == actor.ID {
			// Move the existing actor to the end of the list
			notification.Actors = append(notification.Actors[:i], notification.Actors[i+1:]...) // Remove the actor
			notification.Actors = append(notification.Actors, existingActor)                    // Append the actor to the end
			actorExists = true
			break
		}
	}
	if !actorExists {
		notification.Actors = append(notification.Actors, actor)
	}

	// Add any new action types that are not already present or move existing to the end
	for _, newActionType := range newActionTypes {
		actionExists := false
		for i, existingActionType := range notification.ActionType {
			if existingActionType == newActionType {
				// Move the existing action type to the end of the list
				notification.ActionType = append(notification.ActionType[:i], notification.ActionType[i+1:]...) // Remove the action type
				notification.ActionType = append(notification.ActionType, existingActionType)                   // Append the action type to the end
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
	notification.ReferenceContent = ReferenceContent

	// Save the updated notification
	if err := storage.SaveNotification(notification); err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

// CreateNewNotification creates a new notification entry
func createNewNotification(userID string, actor models.Actor, actionTypes []string, referenceID uuid.UUID, ReferenceContent string) (*models.Notification, error) {
	newNotification := models.Notification{
		ID:                  uuid.New(),
		UserID:              userID,
		Actors:              []models.Actor{actor}, // Adding the actor to the Actors array
		ActionType:          actionTypes,           // Using the array of action types
		ReferenceID:         referenceID,
		NotificationContent: "",
		ReferenceContent:    ReferenceContent,
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
func FetchUserNotificationsService(userID string, limit int, lang string) ([]models.Notification, error) {
	notifications, err := storage.FetchNotificationsByUserID(userID, limit)
	if err != nil {
		log.Println("Error fetching user notifications:", err)
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}

	// Process notifications to create message content
	for i := range notifications {
		notifications[i].NotificationContent = utils.CreateNotificationMessage(notifications[i], lang)
	}

	return notifications, nil
}
