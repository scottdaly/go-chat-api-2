package models

import "gorm.io/gorm"

type Persona struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChatMessage struct {
	gorm.Model
	PersonaID uint   `json:"persona_id"`
	Message   string `json:"message"`
	Response  string `json:"response"`
}