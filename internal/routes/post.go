package routes

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func PostsRoutesSetup(app *fiber.App) {

	app.Post("/posts", func(c *fiber.Ctx) error {
		return handlers.CreatePost(c)
	})
}
