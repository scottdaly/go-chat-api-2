package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rscottdaly/go-chat-api-2/claude"
	"github.com/rscottdaly/go-chat-api-2/database"
	"github.com/rscottdaly/go-chat-api-2/models"
)

func ChatHandler(c *fiber.Ctx) error {
	var req struct {
		ConversationID uint   `json:"conversation_id"`
		PersonaID      uint   `json:"persona_id"`
		Message        string `json:"message"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON: " + err.Error()})
	}

	var conversation models.Conversation
	var persona models.Persona

	if req.ConversationID != 0 {
		if err := database.DB.Preload("Messages").First(&conversation, req.ConversationID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Conversation not found"})
		}
		persona = conversation.Persona
	} else {
		if err := database.DB.First(&persona, req.PersonaID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Persona not found"})
		}
		conversation = models.Conversation{PersonaID: persona.ID}
		database.DB.Create(&conversation)
	}

	userMessage := models.Message{
		ConversationID: conversation.ID,
		Role:           "user",
		Content:        req.Message,
	}
	database.DB.Create(&userMessage)

	conversation.Messages = append(conversation.Messages, userMessage)

	// Generate response using Claude API
	response, err := claude.GenerateResponse(persona, conversation.Messages)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate response: " + err.Error()})
	}

	aiMessage := models.Message{
		ConversationID: conversation.ID,
		Role:           "ai",
		Content:        response,
	}
	database.DB.Create(&aiMessage)

	return c.JSON(fiber.Map{
		"conversation_id": conversation.ID,
		"response":        response,
	})
}

func GetConversationHandler(c *fiber.Ctx) error {
	conversationID := c.Params("id")

	var conversation models.Conversation
	if err := database.DB.Preload("Messages").First(&conversation, conversationID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Conversation not found"})
	}

	return c.JSON(conversation)
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

func ListPersonasHandler(c *fiber.Ctx) error {
	var personas []models.Persona
	database.DB.Find(&personas)
	return c.JSON(personas)
}
