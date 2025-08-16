package dto

import (
	"time"
)

type CreateEventRequest struct {
	Title                   string     `json:"title" binding:"required,max=500"`
	Location                *string    `json:"location,omitempty" binding:"omitempty,max=500"`
	StartsAt                time.Time  `json:"starts_at" binding:"required"`
	RemindBeforeMinutes     *int       `json:"remind_before_minutes,omitempty" binding:"omitempty,min=0,max=10080"`
	RemindFrequencyMinutes  *int       `json:"remind_frequency_minutes,omitempty" binding:"omitempty,min=1,max=1440"`
	RequireConfirmation     *bool      `json:"require_confirmation,omitempty"`
	MaxNotifications        *int       `json:"max_notifications,omitempty" binding:"omitempty,min=1,max=10"`
}

type UpdateEventRequest struct {
	Title                   *string    `json:"title,omitempty" binding:"omitempty,max=500"`
	Location                *string    `json:"location,omitempty" binding:"omitempty,max=500"`
	StartsAt                *time.Time `json:"starts_at,omitempty"`
	RemindBeforeMinutes     *int       `json:"remind_before_minutes,omitempty" binding:"omitempty,min=0,max=10080"`
	RemindFrequencyMinutes  *int       `json:"remind_frequency_minutes,omitempty" binding:"omitempty,min=1,max=1440"`
	RequireConfirmation     *bool      `json:"require_confirmation,omitempty"`
	MaxNotifications        *int       `json:"max_notifications,omitempty" binding:"omitempty,min=1,max=10"`
	Status                  *string    `json:"status,omitempty" binding:"omitempty,oneof=scheduled confirmed canceled completed"`
}

type UpdateUserProfileRequest struct {
	Name                            *string `json:"name,omitempty" binding:"omitempty,max=255"`
	Timezone                        *string `json:"timezone,omitempty" binding:"omitempty,max=50"`
	DefaultRemindBeforeMinutes      *int    `json:"default_remind_before_minutes,omitempty" binding:"omitempty,min=0,max=10080"`
	DefaultRemindFrequencyMinutes   *int    `json:"default_remind_frequency_minutes,omitempty" binding:"omitempty,min=1,max=1440"`
	DefaultRequireConfirmation      *bool   `json:"default_require_confirmation,omitempty"`
	LLMProvider                     *string `json:"llm_provider,omitempty"`
	LLMModel                        *string `json:"llm_model,omitempty"`
}

type ListEventsQuery struct {
	StartDate *time.Time `form:"start_date"`
	EndDate   *time.Time `form:"end_date"`
	Status    *string    `form:"status" binding:"omitempty,oneof=scheduled confirmed canceled completed"`
	Limit     *int       `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset    *int       `form:"offset" binding:"omitempty,min=0"`
}

type AuthenticateRequest struct {
	WANumber string `json:"wa_number" binding:"required,max=20"`
}