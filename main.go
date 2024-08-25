package main

import (
	"log"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.Setup(app)

	log.Fatal(app.Listen(":4000"))
}
