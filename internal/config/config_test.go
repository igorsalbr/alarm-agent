package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWhitelistNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single number",
			input:    "+5511999999999",
			expected: []string{"+5511999999999"},
		},
		{
			name:     "multiple numbers",
			input:    "+5511999999999,+5511888888888,+5511777777777",
			expected: []string{"+5511999999999", "+5511888888888", "+5511777777777"},
		},
		{
			name:     "numbers with spaces",
			input:    " +5511999999999 , +5511888888888 , +5511777777777 ",
			expected: []string{"+5511999999999", "+5511888888888", "+5511777777777"},
		},
		{
			name:     "empty entries",
			input:    "+5511999999999,,+5511888888888",
			expected: []string{"+5511999999999", "+5511888888888"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseWhitelistNumbers(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectedErr string
	}{
		{
			name: "valid config",
			config: Config{
				Infobip: InfobipConfig{
					APIKey:         "test-key",
					WhatsAppSender: "test-sender",
				},
				LLM: LLMConfig{
					Provider:     "anthropic",
					AnthropicKey: "test-anthropic-key",
				},
				Security: SecurityConfig{
					WhitelistNumbers: []string{"+5511999999999"},
				},
			},
			expectedErr: "",
		},
		{
			name: "missing infobip api key",
			config: Config{
				Infobip: InfobipConfig{
					WhatsAppSender: "test-sender",
				},
				LLM: LLMConfig{
					Provider:     "anthropic",
					AnthropicKey: "test-anthropic-key",
				},
				Security: SecurityConfig{
					WhitelistNumbers: []string{"+5511999999999"},
				},
			},
			expectedErr: "INFOBIP_API_KEY is required",
		},
		{
			name: "missing whatsapp sender",
			config: Config{
				Infobip: InfobipConfig{
					APIKey: "test-key",
				},
				LLM: LLMConfig{
					Provider:     "anthropic",
					AnthropicKey: "test-anthropic-key",
				},
				Security: SecurityConfig{
					WhitelistNumbers: []string{"+5511999999999"},
				},
			},
			expectedErr: "INFOBIP_WHATSAPP_SENDER is required",
		},
		{
			name: "invalid llm provider",
			config: Config{
				Infobip: InfobipConfig{
					APIKey:         "test-key",
					WhatsAppSender: "test-sender",
				},
				LLM: LLMConfig{
					Provider: "invalid-provider",
				},
				Security: SecurityConfig{
					WhitelistNumbers: []string{"+5511999999999"},
				},
			},
			expectedErr: "LLM_PROVIDER must be 'anthropic' or 'openai'",
		},
		{
			name: "missing anthropic key when using anthropic",
			config: Config{
				Infobip: InfobipConfig{
					APIKey:         "test-key",
					WhatsAppSender: "test-sender",
				},
				LLM: LLMConfig{
					Provider: "anthropic",
				},
				Security: SecurityConfig{
					WhitelistNumbers: []string{"+5511999999999"},
				},
			},
			expectedErr: "ANTHROPIC_API_KEY is required when using anthropic provider",
		},
		{
			name: "empty whitelist",
			config: Config{
				Infobip: InfobipConfig{
					APIKey:         "test-key",
					WhatsAppSender: "test-sender",
				},
				LLM: LLMConfig{
					Provider:     "anthropic",
					AnthropicKey: "test-anthropic-key",
				},
				Security: SecurityConfig{
					WhitelistNumbers: []string{},
				},
			},
			expectedErr: "WHITELIST_NUMBERS must contain at least one number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	config := &Config{
		App: AppConfig{
			Environment: "development",
		},
	}

	assert.True(t, config.IsDevelopment())
	assert.False(t, config.IsProduction())
}

func TestConfig_IsProduction(t *testing.T) {
	config := &Config{
		App: AppConfig{
			Environment: "production",
		},
	}

	assert.False(t, config.IsDevelopment())
	assert.True(t, config.IsProduction())
}

func TestGetEnvOrDefault(t *testing.T) {
	key := "TEST_ENV_VAR"
	defaultValue := "default"
	testValue := "test"

	// Test default value when env var is not set
	result := getEnvOrDefault(key, defaultValue)
	assert.Equal(t, defaultValue, result)

	// Test env var value when set
	os.Setenv(key, testValue)
	defer os.Unsetenv(key)

	result = getEnvOrDefault(key, defaultValue)
	assert.Equal(t, testValue, result)
}

func TestGetEnvAsIntOrDefault(t *testing.T) {
	key := "TEST_INT_ENV_VAR"
	defaultValue := 42
	testValue := "123"

	// Test default value when env var is not set
	result := getEnvAsIntOrDefault(key, defaultValue)
	assert.Equal(t, defaultValue, result)

	// Test env var value when set to valid int
	os.Setenv(key, testValue)
	defer os.Unsetenv(key)

	result = getEnvAsIntOrDefault(key, defaultValue)
	assert.Equal(t, 123, result)

	// Test default value when env var is set to invalid int
	os.Setenv(key, "invalid")
	result = getEnvAsIntOrDefault(key, defaultValue)
	assert.Equal(t, defaultValue, result)
}
