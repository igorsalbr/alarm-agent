package dto

import (
	"time"

	"github.com/alarm-agent/internal/domain"
)

type EventResponse struct {
	ID                      int        `json:"id"`
	Title                   string     `json:"title"`
	Location                *string    `json:"location,omitempty"`
	StartsAt                time.Time  `json:"starts_at"`
	RemindBeforeMinutes     int        `json:"remind_before_minutes"`
	RemindFrequencyMinutes  int        `json:"remind_frequency_minutes"`
	RequireConfirmation     bool       `json:"require_confirmation"`
	MaxNotifications        int        `json:"max_notifications"`
	Status                  string     `json:"status"`
	NotificationsSent       int        `json:"notifications_sent"`
	LastNotifiedAt          *time.Time `json:"last_notified_at,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

type UserResponse struct {
	ID                              int     `json:"id"`
	WANumber                        string  `json:"wa_number"`
	Name                            *string `json:"name,omitempty"`
	Timezone                        string  `json:"timezone"`
	DefaultRemindBeforeMinutes      int     `json:"default_remind_before_minutes"`
	DefaultRemindFrequencyMinutes   int     `json:"default_remind_frequency_minutes"`
	DefaultRequireConfirmation      bool    `json:"default_require_confirmation"`
	LLMProvider                     *string `json:"llm_provider,omitempty"`
	LLMModel                        *string `json:"llm_model,omitempty"`
	CreatedAt                       time.Time `json:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at"`
}

type LLMProviderResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	Models      []LLMModelResponse `json:"models,omitempty"`
}

type LLMModelResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
	IsActive    bool   `json:"is_active"`
}

type ListEventsResponse struct {
	Events     []EventResponse `json:"events"`
	TotalCount int             `json:"total_count"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func EventToResponse(event *domain.Event) EventResponse {
	return EventResponse{
		ID:                      event.ID,
		Title:                   event.Title,
		Location:                event.Location,
		StartsAt:                event.StartsAt,
		RemindBeforeMinutes:     event.RemindBeforeMinutes,
		RemindFrequencyMinutes:  event.RemindFrequencyMinutes,
		RequireConfirmation:     event.RequireConfirmation,
		MaxNotifications:        event.MaxNotifications,
		Status:                  string(event.Status),
		NotificationsSent:       event.NotificationsSent,
		LastNotifiedAt:          event.LastNotifiedAt,
		CreatedAt:               event.CreatedAt,
		UpdatedAt:               event.UpdatedAt,
	}
}

func UserToResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:                              user.ID,
		WANumber:                        user.WANumber,
		Name:                            user.Name,
		Timezone:                        user.Timezone,
		DefaultRemindBeforeMinutes:      user.DefaultRemindBeforeMinutes,
		DefaultRemindFrequencyMinutes:   user.DefaultRemindFrequencyMinutes,
		DefaultRequireConfirmation:      user.DefaultRequireConfirmation,
		LLMProvider:                     user.LLMProvider,
		LLMModel:                        user.LLMModel,
		CreatedAt:                       user.CreatedAt,
		UpdatedAt:                       user.UpdatedAt,
	}
}

func LLMProviderToResponse(provider *domain.LLMProvider) LLMProviderResponse {
	return LLMProviderResponse{
		ID:          provider.ID,
		Name:        provider.Name,
		DisplayName: provider.DisplayName,
		Description: provider.Description,
		IsActive:    provider.IsActive,
	}
}

func LLMModelToResponse(model *domain.LLMModel) LLMModelResponse {
	return LLMModelResponse{
		ID:          model.ID,
		Name:        model.Name,
		DisplayName: model.DisplayName,
		Description: model.Description,
		IsDefault:   model.IsDefault,
		IsActive:    model.IsActive,
	}
}