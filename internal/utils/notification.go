package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/storage"
)

// MessageTemplates holds the message format for each action type in both Arabic and English
var MessageTemplates = map[string]map[string]string{
	"like": {
		"ar": "أبدى %s إعجاباً بمشاركتك: %s",
		"en": "%s liked your post: %s",
	},
	"comment": {
		"ar": "علق %s على مشاركتك: %s",
		"en": "%s commented on your post: %s",
	},
	"mention": {
		"ar": "أشار إليك %s: %s",
		"en": "%s mentioned you in a comment: %s",
	},
}

// createNotificationMessage generates the notification message based on actions
func CreateNotificationMessage(notification models.Notification, lang string) string {
	if len(notification.Actors) == 0 {
		return ""
	}

	lastActionType, lastActor := CollectLastActionType(notification)

	return BuildNotificationMessage(lastActionType, lastActor, notification, lang)
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
func BuildNotificationMessage(lastActionType, lastActor string, notification models.Notification, lang string) string {
	message := ""

	// Choose the correct message template based on the action type and language
	if template, ok := MessageTemplates[lastActionType]; ok {
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

// SendPushNotification sends a push notification to the user
func SendPushNotification(userID string, notification *models.Notification, lang string) error {
	// Fetch the user's Expo push token from storage
	user, err := storage.FindUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to retrieve push token for user %s: %w", userID, err)
	}

	if user.PushToken == "" {
		return fmt.Errorf("no push token found for user %s", userID)
	}

	// Construct the push notification message
	message := CreateNotificationMessage(*notification, lang)

	// Prepare the payload for the Expo Push Notification API
	payload := map[string]interface{}{
		"to":    user.PushToken,
		"title": notification.Actors[len(notification.Actors)-1].Name, // Customize the title as needed
		"body":  message,
		"data":  map[string]string{"reference_id": notification.ReferenceID.String()}, // Pass the parsed refID
	}

	// Send the notification to the Expo Push Notification service
	response, err := sendToExpoPushAPI(payload)
	if err != nil {
		return fmt.Errorf("failed to send notification to Expo: %w", err)
	}

	// Log the response from the Expo Push Notification service
	log.Println("Response from Expo Push API:", response)

	return nil
}

func sendToExpoPushAPI(payload map[string]interface{}) (string, error) {
	url := "https://exp.host/--/api/v2/push/send"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	return "Notification sent successfully!", nil
}
