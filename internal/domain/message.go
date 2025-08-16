package domain

import (
	"encoding/json"
	"time"
)

type InboundMessage struct {
	ID                int             `json:"id" db:"id"`
	ProviderMessageID string          `json:"provider_message_id" db:"provider_message_id"`
	FromNumber        string          `json:"from_number" db:"from_number"`
	RawPayload        json.RawMessage `json:"raw_payload" db:"raw_payload"`
	ProcessedAt       time.Time       `json:"processed_at" db:"processed_at"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
}

type LLMIntent string

const (
	IntentCreateEvent  LLMIntent = "create_event"
	IntentUpdateEvent  LLMIntent = "update_event"
	IntentCancelEvent  LLMIntent = "cancel_event"
	IntentListEvents   LLMIntent = "list_events"
	IntentConfirmEvent LLMIntent = "confirm_event"
	IntentDeclineEvent LLMIntent = "decline_event"
	IntentSmallTalk    LLMIntent = "small_talk"
	IntentUnknown      LLMIntent = "unknown"
)

type LLMResponse struct {
	Intent            LLMIntent              `json:"intent"`
	Entities          map[string]interface{} `json:"entities"`
	Confidence        float64                `json:"confidence"`
	FollowUpQuestion  *string                `json:"follow_up_question"`
	Notes             *string                `json:"notes"`
}

type EventEntities struct {
	Title                   *string    `json:"title"`
	StartsAt                *time.Time `json:"starts_at"`
	Location                *string    `json:"location"`
	Participants            []string   `json:"participants"`
	RemindBeforeMinutes     *int       `json:"remind_before_minutes"`
	RemindFrequencyMinutes  *int       `json:"remind_frequency_minutes"`
	RequireConfirmation     *bool      `json:"require_confirmation"`
	MaxNotifications        *int       `json:"max_notifications"`
	Identifier              *EventIdentifier `json:"identifier"`
}

type EventIdentifier struct {
	EventID   *int    `json:"event_id"`
	Title     *string `json:"title"`
	DateHint  *string `json:"date_hint"`
}