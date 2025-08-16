package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Infobip  InfobipConfig
	LLM      LLMConfig
	Security SecurityConfig
	Worker   WorkerConfig
}

type AppConfig struct {
	Port            string
	Environment     string
	DefaultTimezone string
}

type DatabaseConfig struct {
	DSN string
}

type InfobipConfig struct {
	BaseURL        string
	APIKey         string
	WhatsAppSender string
	WebhookSecret  string
}

type LLMConfig struct {
	// Keep API keys for backward compatibility during migration
	AnthropicKey string
	OpenAIKey    string
}

type SecurityConfig struct {
	WhitelistNumbers   []string
	RateLimitPerMinute int
}

type WorkerConfig struct {
	ReminderTickInterval time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		App: AppConfig{
			Port:            getEnvOrDefault("PORT", "8080"),
			Environment:     getEnvOrDefault("ENV", "development"),
			DefaultTimezone: getEnvOrDefault("TIMEZONE_DEFAULT", "America/Sao_Paulo"),
		},
		Database: DatabaseConfig{
			DSN: getEnvOrDefault("POSTGRES_DSN", "postgres://alarm_user:alarm_pass@localhost:5432/alarm_agent?sslmode=disable"),
		},
		Infobip: InfobipConfig{
			BaseURL:        getEnvOrDefault("INFOBIP_BASE_URL", "https://api.infobip.com"),
			APIKey:         os.Getenv("INFOBIP_API_KEY"),
			WhatsAppSender: os.Getenv("INFOBIP_WHATSAPP_SENDER"),
			WebhookSecret:  os.Getenv("INFOBIP_WEBHOOK_SECRET"),
		},
		LLM: LLMConfig{
			AnthropicKey: os.Getenv("ANTHROPIC_API_KEY"),
			OpenAIKey:    os.Getenv("OPENAI_API_KEY"),
		},
		Security: SecurityConfig{
			WhitelistNumbers:   parseWhitelistNumbers(os.Getenv("WHITELIST_NUMBERS")),
			RateLimitPerMinute: getEnvAsIntOrDefault("RATE_LIMIT_PER_MINUTE", 30),
		},
		Worker: WorkerConfig{
			ReminderTickInterval: time.Duration(getEnvAsIntOrDefault("REMINDER_TICK_SECONDS", 30)) * time.Second,
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.Infobip.APIKey == "" {
		return fmt.Errorf("INFOBIP_API_KEY is required")
	}

	if c.Infobip.WhatsAppSender == "" {
		return fmt.Errorf("INFOBIP_WHATSAPP_SENDER is required")
	}

	// LLM configuration is now handled by database, no validation needed here

	if len(c.Security.WhitelistNumbers) == 0 {
		return fmt.Errorf("WHITELIST_NUMBERS must contain at least one number")
	}

	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func parseWhitelistNumbers(numbersStr string) []string {
	if numbersStr == "" {
		return []string{}
	}

	numbers := strings.Split(numbersStr, ",")
	result := make([]string, 0, len(numbers))

	for _, number := range numbers {
		if trimmed := strings.TrimSpace(number); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
