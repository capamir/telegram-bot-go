package bot

import (
    "context"
    "fmt"
    
    "github.com/go-telegram/bot"
    "github.com/capamir/telegram-bot-go/internal/usecase"  // ← Add this
)

// Bot wraps the Telegram bot with custom functionality
type Bot struct {
    *bot.Bot
    usecase *usecase.ChatUsecase  // ← Add this field
}

// New creates and initializes a new Telegram bot
func New(ctx context.Context, token string, uc *usecase.ChatUsecase) (*Bot, error) {  // ← Add uc parameter
    if token == "" {
        return nil, fmt.Errorf("telegram token cannot be empty")
    }

    opts := []bot.Option{
        bot.WithDefaultHandler(defaultHandler),
    }

    b, err := bot.New(token, opts...)
    if err != nil {
        return nil, fmt.Errorf("failed to create bot: %w", err)
    }

    return &Bot{
        Bot: b,
        usecase: uc,  // ← Store usecase
    }, nil
}
