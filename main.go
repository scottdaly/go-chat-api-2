package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rscottdaly/go-chat-api-2/database"
	"github.com/rscottdaly/go-chat-api-2/handlers"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	store             *session.Store
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://venturementor.co/api/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	store = session.New()
}

func main() {
	// Initialize database
	database.InitDatabase()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})

	// CORS middleware
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the AI Chatting API!")
	})

	// Auth routes
	app.Get("/login", handlers.HandleGoogleLogin(googleOauthConfig))
	app.Get("/auth/google/callback", handlers.HandleGoogleCallback(googleOauthConfig, store))
	app.Get("/auth/status", handlers.HandleAuthStatus(store))
	app.Post("/logout", handlers.HandleLogout(store))

	// Protected routes
	app.Use(handlers.AuthMiddleware(store))
	app.Post("/chat", handlers.ChatHandler)
	app.Get("/personas", handlers.ListPersonasHandler)
	app.Post("/personas", handlers.CreatePersonaHandler)
	app.Get("/conversations/:id", handlers.GetConversationHandler)

	log.Println("Server starting on port 8080...")
	log.Fatal(app.Listen(":8080"))
}
