package test

import (
	"context"
	"testing"

	"github.com/capamir/telegram-bot-go/internal/ai"
	"github.com/capamir/telegram-bot-go/internal/config"
)

func TestGeminiIntegration(t *testing.T) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	ctx := context.Background()

	// Create AI client
	client, err := ai.NewClient(ctx, cfg.GeminiAPIKey, cfg.GeminiModel)
	if err != nil {
		t.Fatalf("Failed to create AI client: %v", err)
	}

	t.Logf("Testing with model: %s", client.GetModel())

	// Test cases
	tests := []struct {
		name   string
		prompt string
	}{
		{
			name:   "Simple Question",
			prompt: "What is Go programming? Answer in one sentence.",
		},
		{
			name:   "Technical Query",
			prompt: "Explain Telegram bots briefly.",
		},
		{
			name:   "Creative Task",
			prompt: "Tell me an interesting fact about AI.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.GenerateResponse(ctx, tt.prompt)
			
			if err != nil {
				t.Errorf("GenerateResponse() error = %v", err)
				return
			}
			
			if response == "" {
				t.Error("Expected non-empty response")
				return
			}
			
			// Show first 100 chars of response
			preview := response
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			t.Logf("✅ Response: %s", preview)
		})
	}
}

func TestGeminiClientValidation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		apiKey    string
		model     string
		wantError bool
	}{
		{"Valid Config", "test-key-123", "gemini-2.5-flash", false},
		{"Empty API Key", "", "gemini-2.5-flash", true},
		{"Empty Model", "test-key-123", "", true},
		{"Whitespace API Key", "   ", "gemini-2.5-flash", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ai.NewClient(ctx, tt.apiKey, tt.model)
			
			if (err != nil) != tt.wantError {
				t.Errorf("NewClient() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// Run tests with: go test ./test -v
// Run with coverage: go test ./test -v -cover
