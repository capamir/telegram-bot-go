package main

import (
    "context"
    "log"
    "os"
    "os/signal"

    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get bot token from environment
    token := os.Getenv("TELEGRAM_BOT_TOKEN")
    if token == "" {
        log.Fatal("TELEGRAM_BOT_TOKEN is not set in .env file")
    }

    // Create context that listens for interrupt signals
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    // Create bot options with default handler
    opts := []bot.Option{
        bot.WithDefaultHandler(defaultHandler),
    }

    // Initialize the bot
    b, err := bot.New(token, opts...)
    if err != nil {
        log.Fatal("Failed to create bot:", err)
    }

    // Register command handlers
    b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
    b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler)

    log.Println("Bot started successfully! Send /start to begin...")

    // Start the bot (blocking call)
    b.Start(ctx)
}

// Handler for /start command
func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
    // Send welcome message
    _, err := b.SendMessage(ctx, &bot.SendMessageParams{
        ChatID: update.Message.Chat.ID,
        Text:   "👋 Hello! I'm your Telegram bot!\n\nSend me any message and I'll echo it back.\n\nUse /help to see available commands.",
    })
    if err != nil {
        log.Println("Error sending start message:", err)
    }
}

// Handler for /help command
func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
    helpText := `📚 <b>Available Commands:</b>

/start - Start the bot
/help - Show this help message

Just send me any text message and I'll repeat it back to you!`

    _, err := b.SendMessage(ctx, &bot.SendMessageParams{
        ChatID:    update.Message.Chat.ID,
        Text:      helpText,
        ParseMode: models.ParseModeHTML,  
    })
    if err != nil {
        log.Println("Error sending help message:", err)
    }
}

// Default handler for all other messages (echo functionality)
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
    // Only process text messages
    if update.Message == nil || update.Message.Text == "" {
        return
    }

    // Echo the message back
    _, err := b.SendMessage(ctx, &bot.SendMessageParams{
        ChatID: update.Message.Chat.ID,
        Text:   "You said: " + update.Message.Text,
    })
    if err != nil {
        log.Println("Error sending echo message:", err)
    }
}
