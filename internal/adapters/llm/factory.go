package llm

import (
	"fmt"

	"github.com/alarm-agent/internal/ports"
)

func NewLLMClient(provider, apiKey, model string) (ports.LLMClient, error) {
	switch provider {
	case "anthropic":
		if apiKey == "" {
			return nil, fmt.Errorf("anthropic API key is required")
		}
		return NewAnthropicClient(apiKey, model), nil
	case "openai":
		if apiKey == "" {
			return nil, fmt.Errorf("openai API key is required")
		}
		return NewOpenAIClient(apiKey, model), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", provider)
	}
}
