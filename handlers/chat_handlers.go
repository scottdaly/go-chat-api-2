package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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

	now := time.Now()
	if req.ConversationID == 0 {
		// Creating a new conversation
		conversation = models.Conversation{
			PersonaID:     persona.ID,
			StartedAt:     now,
			LastMessageAt: now,
		}
		database.DB.Create(&conversation)
	}

	userMessage := models.Message{
		ConversationID: conversation.ID,
		Role:           "user",
		Content:        req.Message,
		Timestamp:      now,
	}
	database.DB.Create(&userMessage)

	// Update LastMessageAt for the conversation
	database.DB.Model(&conversation).Update("LastMessageAt", now)

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
		Timestamp:      time.Now(),
	}
	database.DB.Create(&aiMessage)

	// Update LastMessageAt again
	database.DB.Model(&conversation).Update("LastMessageAt", aiMessage.Timestamp)

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

	// Set CreatorID and CreatedAt
	currentUserID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	persona.CreatorID = currentUserID
	persona.CreatedAt = time.Now()

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

// Helper function to get current user ID
func getCurrentUserID(c *fiber.Ctx) (uint, error) {
	store := c.Locals("store").(*session.Store)
	sess, err := store.Get(c)
	if err != nil {
		return 0, err
	}

	user := sess.Get("user")
	if user == nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "User not found in session")
	}

	userModel, ok := user.(models.User)
	if !ok {
		return 0, fiber.NewError(fiber.StatusInternalServerError, "Invalid user data in session")
	}

	return userModel.ID, nil
}
