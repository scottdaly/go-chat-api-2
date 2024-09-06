package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string    `json:"username"`
	Email    string    `json:"email" gorm:"unique"`
	CreatedAt time.Time `json:"created_at"`
}

type Persona struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatorID   uint      `json:"creator_id"`
	Creator     User      `json:"creator" gorm:"foreignKey:CreatorID"`
	CreatedAt   time.Time `json:"created_at"`
}

type Conversation struct {
	gorm.Model
	PersonaID      uint      `json:"persona_id"`
	Persona        Persona
	Messages       []Message `json:"messages"`
	StartedAt      time.Time `json:"started_at"`
	LastMessageAt  time.Time `json:"last_message_at"`
}

type Message struct {
	gorm.Model
	ConversationID uint      `json:"conversation_id"`
	Role           string    `json:"role"` // "user" or "ai"
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
}
