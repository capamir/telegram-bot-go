package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Telegram configuration
	TelegramToken string

	// Gemini AI configuration
	GeminiAPIKey string
	GeminiModel  string
}

// Load reads and validates configuration from environment
func Load() (*Config, error) {
	// Load .env file (ignore error if not exists - might use system env vars)
	_ = godotenv.Load()

	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		GeminiAPIKey:  os.Getenv("GEMINI_API_KEY"),
		GeminiModel:   getEnvOrDefault("GEMINI_MODEL", "gemini-2.5-flash"),
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if all required config is present
func (c *Config) Validate() error {
	if c.TelegramToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if c.GeminiAPIKey == "" {
		return fmt.Errorf("GEMINI_API_KEY is required")
	}
	if c.GeminiModel == "" {
		return fmt.Errorf("GEMINI_MODEL cannot be empty")
	}
	return nil
}

// getEnvOrDefault returns environment variable or default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
