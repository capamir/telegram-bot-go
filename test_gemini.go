package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
)

func TestGemini() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        log.Fatal("GEMINI_API_KEY is not set in .env file")
    }

    ctx := context.Background()

    // Create Gemini client
    geminiClient, err := NewGeminiClient(ctx, apiKey)
    if err != nil {
        log.Fatal("Failed to create Gemini client:", err)
    }

    // Test prompts
    testPrompts := []string{
        "What is Go programming language? Answer in one sentence.",
        "Explain Telegram bots in one sentence.",
        "Tell me a fun fact about AI in one sentence.",
    }

    fmt.Println("🤖 Testing Gemini API Integration")
    fmt.Println("Using model: gemini-2.5-flash")
    fmt.Println("==================================================")

    for i, prompt := range testPrompts {
        fmt.Printf("\nTest %d: %s\n", i+1, prompt)
        
        response, err := geminiClient.GenerateResponse(ctx, prompt)
        if err != nil {
            log.Printf("❌ Error: %v\n", err)
            continue
        }

        fmt.Printf("✅ Response: %s\n", response)
        fmt.Println("--------------------------------------------------")
    }

    fmt.Println("\n✅ Gemini API test completed!")
}
