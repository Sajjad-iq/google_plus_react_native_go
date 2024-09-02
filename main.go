package main

import (
	"log"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Import the CORS middleware
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Connect to the database
	database.Connect()

	// Set up the Fiber app
	app := fiber.New()

	// Configure CORS to allow requests from any origin
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow any origin
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Set up the routes
	routes.PostsRoutesSetup(app)
	routes.UsersRoutesSetup(app)

	// Start the server
	log.Fatal(app.Listen(":4000"))
}
