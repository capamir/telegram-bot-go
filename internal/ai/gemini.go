package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

// Client wraps Gemini API client with additional functionality
type Client struct {
	client  *genai.Client
	model   string
	timeout time.Duration
}

// NewClient creates a new Gemini AI client with validation
func NewClient(ctx context.Context, apiKey, model string) (*Client, error) {
	// Validate inputs
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}
	if strings.TrimSpace(model) == "" {
		return nil, fmt.Errorf("model name cannot be empty")
	}

	// Create Gemini client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &Client{
		client:  client,
		model:   model,
		timeout: 30 * time.Second, // Default 30s timeout
	}, nil
}

// GenerateResponse sends a prompt to Gemini and returns the AI response
func (c *Client) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	// Validate prompt
	if strings.TrimSpace(prompt) == "" {
		return "", fmt.Errorf("prompt cannot be empty")
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Call Gemini API
	result, err := c.client.Models.GenerateContent(
		timeoutCtx,
		c.model,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		// Check for timeout
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("request timed out after %v", c.timeout)
		}
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract response text
	response := strings.TrimSpace(result.Text())
	if response == "" {
		return "I couldn't generate a response. Please try again.", nil
	}

	return response, nil
}

// SetTimeout configures the request timeout duration
func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// GetModel returns the current model name
func (c *Client) GetModel() string {
	return c.model
}
