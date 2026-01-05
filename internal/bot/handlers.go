// internal/bot/handlers.go
package bot

import (
	"context"
	"log"
	"strings"

	"github.com/capamir/telegram-bot-go/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// StartHandler handles the /start command
func (b *Bot) StartHandler(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	welcomeText := `👋 Welcome to AI Bot!

I'm an intelligent assistant powered by Google Gemini AI, built with Go.

ℹ️ How to Use This Bot

To chat with me, choose a tone command:

• /ask - Default friendly conversation
• /serious - Professional, detailed answers
• /satirical - Sarcastic, humorous responses
• /friendly - Casual conversation
• /professional - Business language

Example:
/satirical Why is Go better than Python?

Type /help for more information.`

	_, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   welcomeText,
	})
	if err != nil {
		log.Printf("Error sending start message: %v", err)
	}
}

// HelpHandler handles the /help command
func (b *Bot) HelpHandler(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	helpText := `📚 Help & Information

How to Use:
Simply choose a tone and ask your question!

Tone Commands:
• /ask - Friendly conversation (default)
• /serious - Professional answers
• /satirical - Sarcastic humor
• /friendly - Casual chat
• /professional - Business language

Feature Commands:
• /movie-night [genres] - Get movie recommendations

Other Commands:
/start - Welcome message
/help - This help message

Examples:
/ask What is Clean Architecture?
/serious Explain Go interfaces
/satirical Why do developers love Go?
/movie-night action sci-fi thriller

What I Can Do:
• Answer questions on any topic
• Explain complex concepts
• Write creative content
• Provide recommendations
• And much more!

Powered by Google Gemini 2.5 Flash 🤖✨`

	_, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpText,
	})
	if err != nil {
		log.Printf("Error sending help message: %v", err)
	}
}

// defaultHandler handles all non-command messages
func (b *Bot) defaultHandler(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	helpText := `ℹ️ How to Use This Bot

To chat with me, choose a tone command:

• /ask - Default friendly conversation
• /serious - Professional, detailed answers
• /satirical - Sarcastic, humorous responses
• /friendly - Casual conversation
• /professional - Business language

Example:
/satirical Why is Go better than Python?

Type /help for more information.`

	_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpText,
	})
}

// handleAICommand is a shared helper for all AI tone commands
func (b *Bot) handleAICommand(
	ctx context.Context,
	tgBot *bot.Bot,
	update *models.Update,
	tone domain.Tone,
	commandName string,
) {
	// Extract text after command
	text := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, commandName))

	// Validate
	if text == "" {
		_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "🤔 Please write your question!\n\nExample:\n" + commandName + " What is quantum computing?",
		})
		return
	}

	// Create domain message
	domainMsg := &domain.Message{
		ChatID: update.Message.Chat.ID,
		Text:   text,
		Tone:   tone,
	}

	// Use helper to handle loading + response
	b.sendAIResponse(ctx, tgBot, update.Message.Chat.ID, text, func() (*domain.Response, error) {
		return b.usecase.HandleMessage(ctx, domainMsg)
	})
}

// toneCommandHandler handles all tone-based AI commands dynamically
func (b *Bot) toneCommandHandler(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	// Extract the command (e.g., "/ask", "/serious")
	command := strings.Split(update.Message.Text, " ")[0]

	// Map command to tone
	toneMap := map[string]domain.Tone{
		"/ask":          domain.ToneFriendly,
		"/serious":      domain.ToneSerious,
		"/satirical":    domain.ToneSatirical,
		"/friendly":     domain.ToneFriendly,
		"/professional": domain.ToneProfessional,
	}

	tone, exists := toneMap[command]
	if !exists {
		tone = domain.ToneFriendly // fallback
	}

	b.handleAICommand(ctx, tgBot, update, tone, command)
}

// MovieNightHandler handles /movie-night command for movie recommendations
func (b *Bot) MovieNightHandler(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	fullText := update.Message.Text
	command := "/movie-night"

	// Extract genres
	rawGenres := strings.TrimSpace(strings.TrimPrefix(fullText, command))

	// Question text for display
	questionText := rawGenres
	if questionText == "" {
		questionText = "Movie recommendations based on your preferred genres"
	}

	// Use helper to handle loading + response (CONSISTENT WITH OTHER COMMANDS)
	b.sendAIResponse(ctx, tgBot, update.Message.Chat.ID, questionText, func() (*domain.Response, error) {
		return b.usecase.HandleMovieNight(ctx, rawGenres, update.Message.Chat.ID)
	})
}

func (b *Bot) DiaryHandler(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
    text := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/diary"))
    if text == "" {
        _, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
            ChatID: update.Message.Chat.ID,
            Text:   "📝 Please write your diary entry after the command!\nExample: /diary Today I learned about Dependency Injection in Go.",
        })
        return
    }

    domainMsg := &domain.Message{
        ChatID: update.Message.Chat.ID,
        Text:   text,
    }

    b.sendAIResponse(ctx, tgBot, update.Message.Chat.ID, text, func() (*domain.Response, error) {
        return b.usecase.HandleDiary(ctx, domainMsg)
    })
}
