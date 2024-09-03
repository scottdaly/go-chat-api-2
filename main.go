package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rscottdaly/go-chat-api-2/database"
	"github.com/rscottdaly/go-chat-api-2/handlers"
)

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

	app.Post("/chat", handlers.ChatHandler)
	app.Get("/personas", handlers.ListPersonasHandler)
	app.Post("/personas", handlers.CreatePersonaHandler)
	app.Get("/conversations/:id", handlers.GetConversationHandler)

	log.Println("Server starting on port 8080...")
	log.Fatal(app.Listen(":8080"))
}
