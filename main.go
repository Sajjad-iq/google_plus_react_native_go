package main

import (
	"log"
	"os"

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

	jwt_secret := os.Getenv("JWT_SECRET_KEY")

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

	// Routes that don’t require authentication (like auth routes)
	routes.AuthRoutesSetup(app)
	// Serve static files (uploads)
	app.Static("/uploads", "./uploads")

	// Use JWT middleware for protected routes
	// Only apply JWT for routes that need it
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(jwt_secret)},
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
	}))

	// Protected routes
	routes.PostsRoutesSetup(app)
	routes.UsersRoutesSetup(app)

	// Start the server
	log.Fatal(app.Listen(":4000"))
}
