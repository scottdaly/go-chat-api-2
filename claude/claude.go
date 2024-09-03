package claude

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/rscottdaly/go-chat-api-2/models"
)

const claudeAPIEndpoint = "https://api.anthropic.com/v1/messages"

type ClaudeRequest struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	Max_tokens int       `json:"max_tokens"`
	System     string    `json:"system"`
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

func GenerateResponse(persona models.Persona, conversation []models.Message) (string, error) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		return "", errors.New("CLAUDE_API_KEY environment variable not set")
	}

	client := resty.New()

	messages := []Message{}

	for _, msg := range conversation {
		messages = append(messages, Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	request := ClaudeRequest{
		Model:      "claude-3-5-sonnet-20240620",
		Messages:   messages,
		Max_tokens: 2000,
		System: "You are an AI assistant named " + persona.Name + ". " +
			"Your persona is described as: " + persona.Description + ". " +
			"Please respond to the user's messages in character.",
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", apiKey).
		SetHeader("anthropic-version", "2023-06-01").
		SetBody(request).
		Post(claudeAPIEndpoint)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", errors.New("Claude API request failed with status: " + resp.Status())
	}

	var claudeResp ClaudeResponse
	err = json.Unmarshal(resp.Body(), &claudeResp)
	if err != nil {
		return "", err
	}

	if len(claudeResp.Content) > 0 {
		return claudeResp.Content[0].Text, nil
	}

	return "", errors.New("no content in Claude API response")
}
