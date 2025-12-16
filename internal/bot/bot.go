package bot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/capamir/telegram-bot-go/internal/usecase"
)

// Bot wraps the Telegram bot with custom functionality
type Bot struct {
	*bot.Bot
	usecase *usecase.ChatUsecase
}

// New creates and initializes a new Telegram bot
func New(ctx context.Context, token string, uc *usecase.ChatUsecase) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram token cannot be empty")
	}

	customBot := &Bot{
		usecase: uc,
	}

	// Create bot with default handler as method
	opts := []bot.Option{
		bot.WithDefaultHandler(customBot.defaultHandler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	customBot.Bot = b
	return customBot, nil
}

// RegisterHandlers registers all command handlers for the bot
func (b *Bot) RegisterHandlers() {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, b.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, b.HelpHandler)
}
