package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rscottdaly/go-chat-api-2/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	dbPath := filepath.Join("/app/data", "chat_app.db")
	
	// Ensure the directory exists
	err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create database directory: %v", err)
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	DB.AutoMigrate(&models.Persona{}, &models.ChatMessage{})

	// Seed some initial data only if the table is empty
	var count int64
	DB.Model(&models.Persona{}).Count(&count)
	if count == 0 {
		seedData()
	}
}

func seedData() {
	personas := []models.Persona{
		{Name: "Friendly Assistant", Description: "A helpful and friendly AI assistant."},
		{Name: "Tech Guru", Description: "An AI expert in all things technology."},
		{Name: "Creative Writer", Description: "An AI with a flair for creative writing."},
	}
	DB.Create(&personas)
}
