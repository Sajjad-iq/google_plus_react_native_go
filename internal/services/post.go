package services

import (
	"fmt"
	"io"
	"mime/multipart" // Correct import for FileHeader
	"os"
	"path/filepath"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/google/uuid"
)

// ValidatePostForm validates the form data for creating a post
func ValidatePostForm(form *multipart.Form) (*models.Post, error) {
	post := new(models.Post)

	// Validate and assign form values to the post
	if authorIDs, ok := form.Value["author_id"]; ok && len(authorIDs) > 0 {
		post.AuthorID = authorIDs[0]
	} else {
		return nil, fmt.Errorf("author_id is required")
	}

	if authorNames, ok := form.Value["author_name"]; ok && len(authorNames) > 0 {
		post.AuthorName = authorNames[0]
	} else {
		return nil, fmt.Errorf("author_name is required")
	}

	if bodies, ok := form.Value["body"]; ok && len(bodies) > 0 {
		post.Body = bodies[0]
	}

	if authorAvatar, ok := form.Value["author_avatar"]; ok && len(authorAvatar) > 0 {
		post.AuthorAvatar = authorAvatar[0]
	}

	if shareStates, ok := form.Value["share_state"]; ok && len(shareStates) > 0 {
		post.ShareState = shareStates[0]
	} else {
		return nil, fmt.Errorf("share_state is required")
	}

	// Generate a new UUID for the post
	post.ID = uuid.New()

	return post, nil
}

// SaveImage handles saving the uploaded image and returns the path
func SaveImage(file *multipart.FileHeader) (string, error) {
	imageUUID := uuid.New().String()
	imagePath := filepath.Join("uploads", imageUUID+filepath.Ext(file.Filename))

	// Ensure the directory structure exists
	if err := os.MkdirAll(filepath.Dir(imagePath), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Create a destination file
	dst, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the file content to the destination
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return imagePath, nil
}
