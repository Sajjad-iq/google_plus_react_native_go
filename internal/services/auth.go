package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// IsUserExist checks if a user with the given ID exists in the database.
// It returns the user data and true if the user exists, otherwise nil and false.

// CompareUserData compares the fields of two user objects and returns a map of the changes.
// It also returns a boolean indicating whether there are any differences.
func CompareUserData(existingUser, requestUser *models.User) (map[string]interface{}, bool) {
	changes := make(map[string]interface{})

	// Compare fields and add to the changes map if they differ
	if requestUser.Username != existingUser.Username {
		changes["Username"] = requestUser.Username
	}
	if requestUser.ProfileAvatar != existingUser.ProfileAvatar {
		changes["ProfileAvatar"] = requestUser.ProfileAvatar
	}

	return changes, len(changes) > 0
}

// Function to generate JWT token for a user
func GenerateJWTForUser(user models.User, secretKey string) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"status":   user.Status,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token expiration set to 72 hours
	}

	// Create the token using claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

// UpdateUserChanges applies the changes from requestUser to existingUser if there are differences.
// It returns an error if the update fails.
func UpdateUserNameAndAvatar(existingUser, requestUser *models.User) (*models.User, error) {
	// Compare existing and request user data
	changes, hasChanges := CompareUserData(existingUser, requestUser)
	if !hasChanges {
		// No changes to update, return the existing user
		return existingUser, nil
	}

	// Apply changes to the existing user
	for key, value := range changes {
		switch key {
		case "Username":
			existingUser.Username = value.(string)
		case "ProfileAvatar":
			existingUser.ProfileAvatar = value.(string)
		}
	}

	// Save updated user to the database and update the user's posts
	if err := UpdateUserNameAndAvatarForEachBelongingTable(existingUser); err != nil {
		// Return the error if the update fails
		return nil, err
	}

	// Return the updated user
	return existingUser, nil
}

// Function to verify the Google OAuth token
func VerifyGoogleOAuthToken(token string) (map[string]interface{}, error) {
	// Google user info endpoint
	url := "https://www.googleapis.com/userinfo/v2/me"

	// Create a request with the token in the Authorization header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Send the request to Google's user info endpoint
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("invalid token")
	}

	// Parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
