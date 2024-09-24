package routes

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func PostsRoutesSetup(app *fiber.App) {

	app.Post("/create-post", func(c *fiber.Ctx) error {
		return handlers.CreatePost(c)
	})

	app.Get("/posts", func(c *fiber.Ctx) error {
		return handlers.GetPosts(c)
	})

	app.Get("/posts/:id", func(c *fiber.Ctx) error {
		return handlers.GetPostByID(c)
	})

	app.Put("/posts/:id/like", handlers.LikePost)

	app.Put("/posts/:id/comment", handlers.CreateComment)
	app.Get("/posts/comment/:id", handlers.FetchComments)
}
