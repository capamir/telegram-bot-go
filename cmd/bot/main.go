package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/capamir/telegram-bot-go/internal/adapter/ai"
	"github.com/capamir/telegram-bot-go/internal/bot"
	"github.com/capamir/telegram-bot-go/internal/config"
	"github.com/capamir/telegram-bot-go/internal/repository"
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

	// 3. Ensure data directory exists for JSON storage
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		log.Fatalf("❌ Failed to create data directory: %v", err)
	}

	// 4. Initialize Repositories (Storage Layer)
	diaryFilePath := filepath.Join(dataDir, "diary.json")
	diaryRepo := repository.NewJSONDiaryRepository(diaryFilePath)
	log.Println("✅ Diary repository initialized at:", diaryFilePath)

	// 5. Create Gemini client (Infrastructure / Adapter layer)
	geminiClient, err := ai.NewClient(
		ctx,
		cfg.GeminiAPIKey,
		cfg.GeminiModel,
	)
	if err != nil {
		log.Fatalf("❌ Failed to create Gemini client: %v", err)
	}

	// 6. Create use case (Application layer)
	// Now injecting BOTH the AI client and the Diary repository
	chatUsecase := usecase.NewChatUsecase(geminiClient, diaryRepo)

	// 7. Create Telegram bot (Delivery layer)
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

	// 8. Start bot (blocking until context is cancelled)
	b.Start(ctx)

	log.Println("👋 Bot stopped gracefully")
}