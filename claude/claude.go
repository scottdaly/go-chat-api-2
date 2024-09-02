package claude

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

const claudeAPIEndpoint = "https://api.anthropic.com/v1/messages"

type ClaudeRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	System   string    `json:"system"`
	MaxTokens int    `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func GenerateResponse(personaName, personaDescription, userMessage string) (string, error) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		return "", errors.New("CLAUDE_API_KEY environment variable not set")
	}

	client := resty.New()

	request := ClaudeRequest{
		Model: "claude-3-5-sonnet-20240620",
		System: "You are an AI assistant named " + personaName + ". " +
					"Your persona is described as: " + personaDescription + ". " +
					"Please respond to the user's message in character.",
		MaxTokens: 2000,
		Messages: []Message{
			{
				Role:    "user",
				Content: userMessage,
			},
		},
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", apiKey).
		SetHeader("anthropic-version", "2023-06-01").
		SetBody(request).
		Post(claudeAPIEndpoint)

		if err != nil {
			log.Printf("Error making request to Claude API: %v", err)
			return "", err
		}
	
		if resp.StatusCode() != 200 {
			log.Printf("Claude API request failed. Status: %s, Body: %s", resp.Status(), string(resp.Body()))
			return "", fmt.Errorf("Claude API request failed. Status: %s, Body: %s", resp.Status(), string(resp.Body()))
		}
	
		var claudeResp ClaudeResponse
		err = json.Unmarshal(resp.Body(), &claudeResp)
		if err != nil {
			log.Printf("Error unmarshaling Claude API response: %v", err)
			return "", err
		}
	
		if len(claudeResp.Content) > 0 {
			return claudeResp.Content[0].Text, nil
		}
	
		return "", errors.New("no content in Claude API response")
	}