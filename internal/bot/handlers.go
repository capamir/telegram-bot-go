package bot

import (
	"context"
	"log"
	"strings"

	"github.com/capamir/telegram-bot-go/internal/domain"
	"github.com/capamir/telegram-bot-go/internal/utils"
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

Example:
/satirical Why is Go better than Python?

Type /help for more information.`

	_, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   welcomeText,
		// ← Remove ParseMode line
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
• /ask <question> - Friendly conversation (default)
• /serious <question> - Professional answers
• /satirical <question> - Sarcastic humor
• /friendly <question> - Casual chat

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
	// Only process text messages
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	helpText := `ℹ️ How to Use This Bot

To chat with me, choose a tone command:

• /ask - Default friendly conversation
• /serious - Professional, detailed answers
• /satirical - Sarcastic, humorous responses
• /friendly - Casual conversation

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
    // 1. Extract text after command
    text := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, commandName))
    
    // 2. Validate
    if text == "" {
        _, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
            ChatID: update.Message.Chat.ID,
            Text:   "🤔 Please write your question!\n\nExample:\n" + commandName + " What is quantum computing?",
        })
        return
    }
    
    // 3. Create domain message
    domainMsg := &domain.Message{
        ChatID: update.Message.Chat.ID,
        Text:   text,
        Tone:   tone,
    }
    
    // 4. Use helper to handle loading + response
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

	// Extract everything after "/movie-night"
	rawGenres := strings.TrimSpace(strings.TrimPrefix(fullText, command))

	// Call dedicated usecase method
	response, err := b.usecase.HandleMovieNight(ctx, rawGenres, update.Message.Chat.ID)
	if err != nil {
		log.Printf("Error handling movie-night: %v", err)
		_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Sorry, I couldn't fetch movie recommendations. Please try again.",
		})
		return
	}

	// For movie-night, the "question" is the genres or a default label
	questionText := rawGenres
	if questionText == "" {
		questionText = "Movie-night recommendations based on your preferred genres"
	}

	formattedResponse := utils.FormatResponseWithQuestion(questionText, response.Text)

	_, err = tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      formattedResponse,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Printf("Error sending movie-night response: %v", err)
	}
}
