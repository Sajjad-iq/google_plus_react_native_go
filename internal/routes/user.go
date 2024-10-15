package routes

import (
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func UsersRoutesSetup(app *fiber.App) {
	app.Get("/user/:id", handlers.GetTheUser)
	app.Post("/test", func(c *fiber.Ctx) error { return handlers.OAuthUserLogin(c) })
	app.Get("/notifications", handlers.FetchNotificationsHandler)
	app.Put("/notifications/read/:id", handlers.MarkNotificationsAsReadHandler)
	app.Put("/push-token", handlers.UpdatePushTokenHandler)
}
