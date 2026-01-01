package config

import "testing"

func TestValidateMissingValues(t *testing.T) {
	cfg := &Config{}

	if err := cfg.Validate(); err == nil || err.Error() != "TELEGRAM_BOT_TOKEN is required" {
		t.Fatalf("expected TELEGRAM_BOT_TOKEN error, got %v", err)
	}

	cfg.TelegramToken = "token"
	if err := cfg.Validate(); err == nil || err.Error() != "GEMINI_API_KEY is required" {
		t.Fatalf("expected GEMINI_API_KEY error, got %v", err)
	}

	cfg.GeminiAPIKey = "key"
	if err := cfg.Validate(); err == nil || err.Error() != "GEMINI_MODEL cannot be empty" {
		t.Fatalf("expected GEMINI_MODEL error, got %v", err)
	}
}

func TestLoadDefaultsModel(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("GEMINI_API_KEY", "key")
	t.Setenv("GEMINI_MODEL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	if cfg.GeminiModel != "gemini-2.5-flash" {
		t.Fatalf("expected default model, got %q", cfg.GeminiModel)
	}
}

func TestLoadMissingToken(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "")
	t.Setenv("GEMINI_API_KEY", "key")
	t.Setenv("GEMINI_MODEL", "model")

	_, err := Load()
	if err == nil || err.Error() != "TELEGRAM_BOT_TOKEN is required" {
		t.Fatalf("expected TELEGRAM_BOT_TOKEN error, got %v", err)
	}
}

func TestLoadMissingAPIKey(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GEMINI_MODEL", "model")

	_, err := Load()
	if err == nil || err.Error() != "GEMINI_API_KEY is required" {
		t.Fatalf("expected GEMINI_API_KEY error, got %v", err)
	}
}
