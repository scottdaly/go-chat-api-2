package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rscottdaly/go-chat-api-2/claude"
	"github.com/rscottdaly/go-chat-api-2/database"
	"github.com/rscottdaly/go-chat-api-2/models"
)


func ChatHandler(c *fiber.Ctx) error {
	var req models.ChatMessage
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON: " + err.Error()})
	}

	var persona models.Persona
	result := database.DB.First(&persona, req.PersonaID)
	if result.Error != nil {
		log.Printf("Error finding persona: %v", result.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Persona not found"})
	}

	log.Printf("Generating response for Persona: %s, Description: %s, Message: %s", persona.Name, persona.Description, req.Message)

	// Generate response using Claude API
	response, err := claude.GenerateResponse(persona.Name, persona.Description, req.Message)
	if err != nil {
		log.Printf("Error generating response: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate response: " + err.Error()})
	}

	req.Response = response

	// Save the chat message to the database
	database.DB.Create(&req)

	return c.JSON(req)
}

func ListPersonasHandler(c *fiber.Ctx) error {
	var personas []models.Persona
	database.DB.Find(&personas)
	return c.JSON(personas)
}

func CreatePersonaHandler(c *fiber.Ctx) error {
	var persona models.Persona
	if err := c.BodyParser(&persona); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON: " + err.Error()})
	}

	result := database.DB.Create(&persona)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create persona"})
	}

	return c.Status(fiber.StatusCreated).JSON(persona)
}