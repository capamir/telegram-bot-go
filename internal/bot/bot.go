package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/capamir/telegram-bot-go/internal/usecase"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
    // Command handlers
    b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, b.StartHandler)
    b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, b.HelpHandler)
    
    // All tone commands use the same handler
    b.RegisterHandler(bot.HandlerTypeMessageText, "/ask", bot.MatchTypePrefix, b.toneCommandHandler)
    b.RegisterHandler(bot.HandlerTypeMessageText, "/serious", bot.MatchTypePrefix, b.toneCommandHandler)
    b.RegisterHandler(bot.HandlerTypeMessageText, "/satirical", bot.MatchTypePrefix, b.toneCommandHandler)
    b.RegisterHandler(bot.HandlerTypeMessageText, "/friendly", bot.MatchTypePrefix, b.toneCommandHandler)
    b.RegisterHandler(bot.HandlerTypeMessageText, "/professional", bot.MatchTypePrefix, b.toneCommandHandler)
    
	// Feature commands
	b.RegisterHandler(bot.HandlerTypeMessageText, "/movie-night", bot.MatchTypePrefix, b.MovieNightHandler)
    
	// Default handler
    b.RegisterHandlerMatchFunc(func(update *models.Update) bool {
        return update.Message != nil && update.Message.Text != "" && !strings.HasPrefix(update.Message.Text, "/")
    }, b.defaultHandler)
}
