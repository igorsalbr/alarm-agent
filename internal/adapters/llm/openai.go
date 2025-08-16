package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"

	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type OpenAIClient struct {
	client *openai.Client
	model  string
}

func NewOpenAIClient(apiKey, model string) ports.LLMClient {
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

func (c *OpenAIClient) Chat(ctx context.Context, systemPrompt, userMessage string) (*domain.LLMResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	response, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       c.model,
		Temperature: 0.1,
		MaxTokens:   1024,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userMessage,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("openai API error: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("empty response from openai")
	}

	content := response.Choices[0].Message.Content
	
	var llmResponse domain.LLMResponse
	if err := json.Unmarshal([]byte(content), &llmResponse); err != nil {
		return &domain.LLMResponse{
			Intent:     domain.IntentUnknown,
			Confidence: 0.0,
			Notes:      &content,
		}, nil
	}

	return &llmResponse, nil
}