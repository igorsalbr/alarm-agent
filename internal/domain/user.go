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
	CreatedAt                     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at" db:"updated_at"`
}

type WhitelistNumber struct {
	Number    string    `json:"number" db:"number"`
	Note      *string   `json:"note,omitempty" db:"note"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
