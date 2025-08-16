package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthropic-ai/anthropic-go"

	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type AnthropicClient struct {
	client *anthropic.Client
	model  string
}

func NewAnthropicClient(apiKey, model string) ports.LLMClient {
	return &AnthropicClient{
		client: anthropic.NewClient(anthropic.WithAPIKey(apiKey)),
		model:  model,
	}
}

func (c *AnthropicClient) Chat(ctx context.Context, systemPrompt, userMessage string) (*domain.LLMResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	message := anthropic.MessageRequest{
		Model: c.model,
		MaxTokens: 1024,
		System: []anthropic.SystemMessage{{
			Text: systemPrompt,
		}},
		Messages: []anthropic.Message{{
			Role: anthropic.RoleUser,
			Content: []anthropic.MessageContent{{
				Type: "text",
				Text: userMessage,
			}},
		}},
		Temperature: anthropic.Float(0.1),
	}

	response, err := c.client.Messages.New(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("anthropic API error: %w", err)
	}

	if len(response.Content) == 0 {
		return nil, fmt.Errorf("empty response from anthropic")
	}

	content := response.Content[0].Text
	
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