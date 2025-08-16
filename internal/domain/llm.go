package domain

import (
	"encoding/json"
	"time"
)

type LLMProvider struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type LLMModel struct {
	ID          int             `json:"id" db:"id"`
	ProviderID  int             `json:"provider_id" db:"provider_id"`
	Name        string          `json:"name" db:"name"`
	DisplayName string          `json:"display_name" db:"display_name"`
	Description string          `json:"description" db:"description"`
	IsActive    bool            `json:"is_active" db:"is_active"`
	IsDefault   bool            `json:"is_default" db:"is_default"`
	Config      json.RawMessage `json:"config" db:"config"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	
	Provider *LLMProvider `json:"provider,omitempty"`
}

func (m *LLMModel) GetProviderName() string {
	if m.Provider != nil {
		return m.Provider.Name
	}
	return ""
}