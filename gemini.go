package main

import (
    "context"
    "fmt"

    "google.golang.org/genai"
)

type GeminiClient struct {
    client *genai.Client
    model  string
}

func NewGeminiClient(ctx context.Context, apiKey string) (*GeminiClient, error) {
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey: apiKey,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create Gemini client: %w", err)
    }

    return &GeminiClient{
        client: client,
        model:  "gemini-2.5-flash", // Use this model name from your Python code
    }, nil
}

func (gc *GeminiClient) GenerateResponse(ctx context.Context, prompt string) (string, error) {
    result, err := gc.client.Models.GenerateContent(
        ctx,
        gc.model,
        genai.Text(prompt),
        nil,
    )
    if err != nil {
        return "", fmt.Errorf("failed to generate content: %w", err)
    }

    response := result.Text()
    if response == "" {
        return "Sorry, I couldn't generate a response.", nil
    }

    return response, nil
}
