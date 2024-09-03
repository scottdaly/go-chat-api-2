package models

import "gorm.io/gorm"

type Persona struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Conversation struct {
	gorm.Model
	PersonaID uint `json:"persona_id"`
	Persona   Persona
	Messages  []Message `json:"messages"`
}

type Message struct {
	gorm.Model
	ConversationID uint   `json:"conversation_id"`
	Role           string `json:"role"` // "user" or "ai"
	Content        string `json:"content"`
}
