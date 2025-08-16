package ports

import (
	"context"
	"time"

	"github.com/alarm-agent/internal/domain"
)

type LLMClient interface {
	Chat(ctx context.Context, systemPrompt, userMessage string) (*domain.LLMResponse, error)
}

type WhatsAppSender interface {
	SendText(ctx context.Context, to, text string) error
}

type WhatsAppWebhookVerifier interface {
	VerifySignature(payload []byte, signature string) bool
}

type TimeProvider interface {
	Now() time.Time
	Sleep(duration time.Duration)
}
