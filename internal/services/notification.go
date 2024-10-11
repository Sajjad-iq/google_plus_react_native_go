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
func CreateOrUpdateNotification(userID, actorID string, actionTypes []string, referenceID uuid.UUID, ReferenceContent string) (*models.Notification, error) {
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
		if err := UpdateExistingNotification(existingNotification, actor, actionTypes, ReferenceContent); err != nil {
			log.Println("Error updating existing notification:", err)
			return nil, err
		}
		return existingNotification, nil
	}

	// Create a new notification
	return CreateNewNotification(userID, actor, actionTypes, referenceID, ReferenceContent)
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
func UpdateExistingNotification(notification *models.Notification, actor models.Actor, newActionTypes []string, ReferenceContent string) error {
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
	notification.ReferenceContent = ReferenceContent

	// Save the updated notification
	if err := storage.SaveNotification(notification); err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

// CreateNewNotification creates a new notification entry
func CreateNewNotification(userID string, actor models.Actor, actionTypes []string, referenceID uuid.UUID, ReferenceContent string) (*models.Notification, error) {
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
		notifications[i].NotificationContent = createNotificationMessage(notifications[i], lang)
	}

	return notifications, nil
}

// createNotificationMessage generates the notification message based on actions
func createNotificationMessage(notification models.Notification, lang string) string {
	if len(notification.Actors) == 0 {
		return ""
	}

	lastActionType, lastActor := CollectLastActionType(notification)

	return buildMessage(lastActionType, lastActor, notification, lang)
}

// CollectLastActionType collects the last action type and the last actor separately
func CollectLastActionType(notification models.Notification) (string, string) {
	var lastActionType, lastActor string

	// Collect the last action type
	if len(notification.ActionType) > 0 {
		lastActionType = notification.ActionType[len(notification.ActionType)-1]
	}

	// Collect the last actor
	if len(notification.Actors) > 0 {
		lastActor = notification.Actors[len(notification.Actors)-1].Name
	}

	return lastActionType, lastActor
}

// buildMessage constructs the final message based on the action type, actor, and language
func buildMessage(lastActionType, lastActor string, notification models.Notification, lang string) string {
	message := ""

	// Choose the correct message template based on the action type and language
	if template, ok := utils.MessageTemplates[lastActionType]; ok {
		// Add "و آخرون" or "and others" if there are multiple actors
		if len(notification.Actors) > 1 {
			if lang == "ar" {
				message += fmt.Sprintf("\u200F%s و آخرون", lastActor)
			} else {
				message += fmt.Sprintf("%s and others", lastActor)
			}
		} else {
			message += fmt.Sprintf("\u200F%s", lastActor)
		}

		// Append the main message
		message = fmt.Sprintf(template[lang], message, notification.ReferenceContent)
	}

	return message
}
