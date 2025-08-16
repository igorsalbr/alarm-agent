package dto

// UpdateUserConfigRequest represents a request to update user configuration
type UpdateUserConfigRequest struct {
	Name                          *string `json:"name,omitempty"`
	Timezone                      *string `json:"timezone,omitempty"`
	DefaultRemindBeforeMinutes    *int    `json:"default_remind_before_minutes,omitempty"`
	DefaultRemindFrequencyMinutes *int    `json:"default_remind_frequency_minutes,omitempty"`
	DefaultRequireConfirmation    *bool   `json:"default_require_confirmation,omitempty"`
	LLMProvider                   *string `json:"llm_provider,omitempty"`
	LLMModel                      *string `json:"llm_model,omitempty"`
	RateLimitPerMinute            *int    `json:"rate_limit_per_minute,omitempty"`
	IsActive                      *bool   `json:"is_active,omitempty"`
}

// AddAllowedContactRequest represents a request to add an allowed contact
type AddAllowedContactRequest struct {
	ContactNumber string  `json:"contact_number" binding:"required"`
	Note          *string `json:"note,omitempty"`
}

// UpdateAllowedContactRequest represents a request to update an allowed contact
type UpdateAllowedContactRequest struct {
	Note *string `json:"note,omitempty"`
}

// UserConfigResponse represents the user's configuration
type UserConfigResponse struct {
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

// AllowedContactResponse represents an allowed contact
type AllowedContactResponse struct {
	ID            int     `json:"id"`
	ContactNumber string  `json:"contact_number"`
	Note          *string `json:"note,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
