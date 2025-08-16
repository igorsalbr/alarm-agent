package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type EventUseCase struct {
	repos ports.Repositories
}

func NewEventUseCase(repos ports.Repositories) *EventUseCase {
	return &EventUseCase{repos: repos}
}

func (uc *EventUseCase) CreateEvent(ctx context.Context, userID int, entities *domain.EventEntities) (*domain.Event, error) {
	if entities.Title == nil || *entities.Title == "" {
		return nil, fmt.Errorf("event title is required")
	}

	if entities.StartsAt == nil {
		return nil, fmt.Errorf("event start time is required")
	}

	user, err := uc.getUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	event := &domain.Event{
		UserID:   userID,
		Title:    *entities.Title,
		StartsAt: *entities.StartsAt,
		Status:   domain.EventStatusScheduled,
	}

	if entities.Location != nil {
		event.Location = entities.Location
	}

	event.RemindBeforeMinutes = getIntOrDefault(entities.RemindBeforeMinutes, user.DefaultRemindBeforeMinutes)
	event.RemindFrequencyMinutes = getIntOrDefault(entities.RemindFrequencyMinutes, user.DefaultRemindFrequencyMinutes)
	event.RequireConfirmation = getBoolOrDefault(entities.RequireConfirmation, user.DefaultRequireConfirmation)
	event.MaxNotifications = getIntOrDefault(entities.MaxNotifications, 3)

	if err := uc.repos.Event().Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

func (uc *EventUseCase) UpdateEvent(ctx context.Context, userID int, entities *domain.EventEntities) (*domain.Event, error) {
	if entities.Identifier == nil {
		return nil, fmt.Errorf("event identifier is required for update")
	}

	events, err := uc.repos.Event().FindByUserAndIdentifier(ctx, userID, entities.Identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find events: %w", err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("event not found")
	}

	if len(events) > 1 {
		return nil, fmt.Errorf("multiple events found, please be more specific")
	}

	event := &events[0]

	if entities.Title != nil {
		event.Title = *entities.Title
	}
	if entities.StartsAt != nil {
		event.StartsAt = *entities.StartsAt
	}
	if entities.Location != nil {
		event.Location = entities.Location
	}
	if entities.RemindBeforeMinutes != nil {
		event.RemindBeforeMinutes = *entities.RemindBeforeMinutes
	}
	if entities.RemindFrequencyMinutes != nil {
		event.RemindFrequencyMinutes = *entities.RemindFrequencyMinutes
	}
	if entities.RequireConfirmation != nil {
		event.RequireConfirmation = *entities.RequireConfirmation
	}
	if entities.MaxNotifications != nil {
		event.MaxNotifications = *entities.MaxNotifications
	}

	if err := uc.repos.Event().Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return event, nil
}

func (uc *EventUseCase) CancelEvent(ctx context.Context, userID int, identifier *domain.EventIdentifier) (*domain.Event, error) {
	events, err := uc.repos.Event().FindByUserAndIdentifier(ctx, userID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find events: %w", err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("event not found")
	}

	if len(events) > 1 {
		return nil, fmt.Errorf("multiple events found, please be more specific")
	}

	event := &events[0]
	event.Status = domain.EventStatusCanceled

	if err := uc.repos.Event().Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to cancel event: %w", err)
	}

	return event, nil
}

func (uc *EventUseCase) ListEvents(ctx context.Context, userID int, startDate, endDate *time.Time) ([]domain.Event, error) {
	if startDate != nil && endDate != nil {
		return uc.repos.Event().GetByUserIDAndDateRange(ctx, userID, *startDate, *endDate)
	}

	events, err := uc.repos.Event().GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	now := time.Now()
	var filteredEvents []domain.Event
	for _, event := range events {
		if event.Status != domain.EventStatusCanceled && event.StartsAt.After(now) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents, nil
}

func (uc *EventUseCase) ConfirmEvent(ctx context.Context, userID int, identifier *domain.EventIdentifier) (*domain.Event, error) {
	events, err := uc.repos.Event().FindByUserAndIdentifier(ctx, userID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find events: %w", err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("event not found")
	}

	if len(events) > 1 {
		return nil, fmt.Errorf("multiple events found, please be more specific")
	}

	event := &events[0]
	if event.Status == domain.EventStatusScheduled {
		event.Status = domain.EventStatusConfirmed
		if err := uc.repos.Event().Update(ctx, event); err != nil {
			return nil, fmt.Errorf("failed to confirm event: %w", err)
		}
	}

	return event, nil
}

func (uc *EventUseCase) GetEventByID(ctx context.Context, userID, eventID int) (*domain.Event, error) {
	event, err := uc.repos.Event().GetByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if event == nil || event.UserID != userID {
		return nil, fmt.Errorf("event not found or access denied")
	}

	return event, nil
}

func (uc *EventUseCase) getUserByID(ctx context.Context, userID int) (*domain.User, error) {
	// This is a placeholder - we need to implement GetByID in UserRepository
	// For now, we'll need to work around this limitation
	return &domain.User{
		ID:                              userID,
		DefaultRemindBeforeMinutes:      30,
		DefaultRemindFrequencyMinutes:   15,
		DefaultRequireConfirmation:      true,
	}, nil
}

func getIntOrDefault(value *int, defaultValue int) int {
	if value != nil {
		return *value
	}
	return defaultValue
}

func getBoolOrDefault(value *bool, defaultValue bool) bool {
	if value != nil {
		return *value
	}
	return defaultValue
}