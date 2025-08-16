package ports

import (
	"context"
	"time"

	"github.com/alarm-agent/internal/domain"
)

type UserRepository interface {
	GetByWANumber(ctx context.Context, waNumber string) (*domain.User, error)
	GetByID(ctx context.Context, userID int) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	UpdateConfig(ctx context.Context, userID int, config *domain.UserConfig) error
}

type WhitelistRepository interface {
	IsWhitelisted(ctx context.Context, number string) (bool, error)
	Add(ctx context.Context, whitelist *domain.WhitelistNumber) error
	Remove(ctx context.Context, number string) error
	List(ctx context.Context) ([]domain.WhitelistNumber, error)
}

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*domain.Event, error)
	GetByUserID(ctx context.Context, userID int) ([]domain.Event, error)
	GetByUserIDAndDateRange(ctx context.Context, userID int, start, end time.Time) ([]domain.Event, error)
	GetPendingReminders(ctx context.Context, reminderWindow time.Duration) ([]domain.EventWithUser, error)
	FindByUserAndIdentifier(ctx context.Context, userID int, identifier *domain.EventIdentifier) ([]domain.Event, error)
}

type InboundMessageRepository interface {
	Create(ctx context.Context, message *domain.InboundMessage) error
	Exists(ctx context.Context, providerMessageID string) (bool, error)
}

type LLMConfigRepository interface {
	GetDefaultModel(ctx context.Context) (*domain.LLMModel, error)
	GetModelByProviderAndName(ctx context.Context, provider, model string) (*domain.LLMModel, error)
	GetActiveProviders(ctx context.Context) ([]domain.LLMProvider, error)
	GetActiveModelsByProvider(ctx context.Context, providerID int) ([]domain.LLMModel, error)
	GetUserLLMConfig(ctx context.Context, userID int) (*domain.LLMModel, error)
}

type UserAllowedContactRepository interface {
	IsAllowed(ctx context.Context, userID int, contactNumber string) (bool, error)
	Add(ctx context.Context, contact *domain.UserAllowedContact) error
	Remove(ctx context.Context, userID int, contactNumber string) error
	List(ctx context.Context, userID int) ([]domain.UserAllowedContact, error)
	GetByUserAndNumber(ctx context.Context, userID int, contactNumber string) (*domain.UserAllowedContact, error)
}

type Repositories interface {
	User() UserRepository
	Whitelist() WhitelistRepository
	Event() EventRepository
	InboundMessage() InboundMessageRepository
	LLMConfig() LLMConfigRepository
	UserAllowedContact() UserAllowedContactRepository
	WithTx(ctx context.Context, fn func(Repositories) error) error
}
