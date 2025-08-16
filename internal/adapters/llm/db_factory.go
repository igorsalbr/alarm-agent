package llm

import (
	"context"
	"fmt"

	"github.com/alarm-agent/internal/config"
	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

func NewLLMClientFromDB(ctx context.Context, repo ports.LLMConfigRepository, cfg *config.Config, userID int) (ports.LLMClient, error) {
	model, err := repo.GetUserLLMConfig(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM configuration: %w", err)
	}

	if model.Provider == nil {
		return nil, fmt.Errorf("model provider information is missing")
	}

	// Get API key from environment based on provider
	var apiKey string
	switch model.Provider.Name {
	case "anthropic":
		apiKey = cfg.LLM.AnthropicKey
	case "openai":
		apiKey = cfg.LLM.OpenAIKey
	default:
		return nil, fmt.Errorf("unsupported provider: %s", model.Provider.Name)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key not found for provider %s", model.Provider.Name)
	}

	return NewLLMClient(model.Provider.Name, apiKey, model.Name)
}