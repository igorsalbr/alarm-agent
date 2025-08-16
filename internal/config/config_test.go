package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// parseWhitelistNumbers function removed - whitelist is now user-level

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
					AnthropicKey: "test-anthropic-key",
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
					AnthropicKey: "test-anthropic-key",
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
					AnthropicKey: "test-anthropic-key",
				},
			},
			expectedErr: "INFOBIP_WHATSAPP_SENDER is required",
		},
		{
			name: "no llm validation needed",
			config: Config{
				Infobip: InfobipConfig{
					APIKey:         "test-key",
					WhatsAppSender: "test-sender",
				},
				LLM: LLMConfig{
					AnthropicKey: "",
					OpenAIKey:    "",
				},
			},
			expectedErr: "",
		},
		{
			name: "missing anthropic key is ok - handled by database",
			config: Config{
				Infobip: InfobipConfig{
					APIKey:         "test-key",
					WhatsAppSender: "test-sender",
				},
				LLM: LLMConfig{
				},
			},
			expectedErr: "",
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
	_ = os.Setenv(key, testValue)
	defer func() { _ = os.Unsetenv(key) }()

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
	_ = os.Setenv(key, testValue)
	defer func() { _ = os.Unsetenv(key) }()

	result = getEnvAsIntOrDefault(key, defaultValue)
	assert.Equal(t, 123, result)

	// Test default value when env var is set to invalid int
	_ = os.Setenv(key, "invalid")
	result = getEnvAsIntOrDefault(key, defaultValue)
	assert.Equal(t, defaultValue, result)
}
