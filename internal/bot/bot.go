package bot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

// Bot wraps the Telegram bot with custom functionality
type Bot struct {
	*bot.Bot
}

// New creates and initializes a new Telegram bot
func New(ctx context.Context, token string) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram token cannot be empty")
	}

	// Create bot with default handler
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &Bot{Bot: b}, nil
}

// RegisterHandlers registers all command handlers for the bot
func (b *Bot) RegisterHandlers() {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, HelpHandler)
}
