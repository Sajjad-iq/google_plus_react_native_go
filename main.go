package main

import (
	"log"

	"github.com/Sajjad-iq/google_plus_react_native_go/internal/database"
	"github.com/Sajjad-iq/google_plus_react_native_go/internal/routes"
	"github.com/Sajjad-iq/google_plus_react_native_go/middleware"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	database.Connect()

	// Set up the Fiber app
	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Register the Logger middleware
	app.Use(middleware.Logger)

	// Set up the routes
	routes.PostsRoutesSetup(app)
	routes.UsersRoutesSetup(app)
	routes.AuthRoutesSetup(app)
	app.Static("/uploads", "./uploads")

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))
	// Start the server
	log.Fatal(app.Listen(":4000"))
}
