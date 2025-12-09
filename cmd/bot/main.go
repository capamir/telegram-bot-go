package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/capamir/telegram-bot-go/internal/bot"
	"github.com/capamir/telegram-bot-go/internal/config"
)

func main() {
	log.Println("🚀 Starting Telegram Bot...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Configuration error: %v", err)
	}
	log.Println("✅ Configuration loaded")

	// Create context with interrupt signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Initialize bot
	b, err := bot.New(ctx, cfg.TelegramToken)
	if err != nil {
		log.Fatalf("❌ Failed to create bot: %v", err)
	}
	log.Println("✅ Bot initialized")

	// Register all command handlers
	b.RegisterHandlers()
	log.Println("✅ Handlers registered")

	// Start bot
	log.Println("✅ Bot is running! Press Ctrl+C to stop.")
	log.Println("📱 Send /start to your bot to begin")
	
	b.Start(ctx)

	log.Println("👋 Bot stopped gracefully")
}
