package main

import (
    "context"
    "log"
    "os"
    "os/signal"

    "github.com/capamir/telegram-bot-go/internal/adapter/ai"
    "github.com/capamir/telegram-bot-go/internal/bot"
    "github.com/capamir/telegram-bot-go/internal/config"
    "github.com/capamir/telegram-bot-go/internal/usecase"
)

func main() {
    log.Println("🚀 Starting Telegram Bot...")

    // 1. Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("❌ Configuration error: %v", err)
    }
    log.Println("✅ Configuration loaded")

    // 2. Root context with graceful shutdown
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    // 3. Create Gemini client (infrastructure / adapter layer)
    geminiClient, err := ai.NewClient(
        ctx,
        cfg.GeminiAPIKey,
        cfg.GeminiModel,
    )
    if err != nil {
        log.Fatalf("❌ Failed to create Gemini client: %v", err)
    }

    // 4. Create use case (application layer)
    chatUsecase := usecase.NewChatUsecase(geminiClient)

    // 5. Create Telegram bot (delivery layer)
    b, err := bot.New(
        ctx,
        cfg.TelegramToken,
        chatUsecase,
    )
    if err != nil {
        log.Fatalf("❌ Failed to start bot: %v", err)
    }

	b.RegisterHandlers()
	log.Println("✅ Handlers registered")
	
    // 6. Start bot (blocking until context is cancelled)
    b.Start(ctx)

    log.Println("👋 Bot stopped gracefully")
}
