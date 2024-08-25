package handlers

import (
	"fmt"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/models"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/services"
	"github.com/gofiber/fiber/v2"
)

func CreatePost(c *fiber.Ctx) error {
	var post models.Post
	fmt.Println("Raw JSON body:", string(c.Body()))

	// Parse the JSON body into the Post struct
	if err := c.BodyParser(&post); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	// Set default values if necessary (e.g., for Likes or CommentsCount)
	post.Likes = 0
	post.CommentsCount = 0

	// Call the service to create the post (no return value expected)
	services.CreatePost(&post)

	return c.Status(fiber.StatusCreated).JSON(post)
}

func GetPosts(c *fiber.Ctx) error {
	posts := services.GetAllPosts()
	return c.JSON(posts)
}
