package domain

import (
	"time"
)

type User struct {
	ID                            int       `json:"id" db:"id"`
	WANumber                      string    `json:"wa_number" db:"wa_number"`
	Name                          *string   `json:"name,omitempty" db:"name"`
	Timezone                      string    `json:"timezone" db:"timezone"`
	DefaultRemindBeforeMinutes    int       `json:"default_remind_before_minutes" db:"default_remind_before_minutes"`
	DefaultRemindFrequencyMinutes int       `json:"default_remind_frequency_minutes" db:"default_remind_frequency_minutes"`
	DefaultRequireConfirmation    bool      `json:"default_require_confirmation" db:"default_require_confirmation"`
	LLMProvider                   *string   `json:"llm_provider,omitempty" db:"llm_provider"`
	LLMModel                      *string   `json:"llm_model,omitempty" db:"llm_model"`
	RateLimitPerMinute            int       `json:"rate_limit_per_minute" db:"rate_limit_per_minute"`
	IsActive                      bool      `json:"is_active" db:"is_active"`
	CreatedAt                     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at" db:"updated_at"`
}

type WhitelistNumber struct {
	Number    string    `json:"number" db:"number"`
	Note      *string   `json:"note,omitempty" db:"note"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserAllowedContact struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	ContactNumber string    `json:"contact_number" db:"contact_number"`
	Note          *string   `json:"note,omitempty" db:"note"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type UserConfig struct {
	UserID                        int     `json:"user_id"`
	Name                          *string `json:"name,omitempty"`
	Timezone                      string  `json:"timezone"`
	DefaultRemindBeforeMinutes    int     `json:"default_remind_before_minutes"`
	DefaultRemindFrequencyMinutes int     `json:"default_remind_frequency_minutes"`
	DefaultRequireConfirmation    bool    `json:"default_require_confirmation"`
	LLMProvider                   *string `json:"llm_provider,omitempty"`
	LLMModel                      *string `json:"llm_model,omitempty"`
	RateLimitPerMinute            int     `json:"rate_limit_per_minute"`
	IsActive                      bool    `json:"is_active"`
}
