package bot

import (
	"context"
	"log"

	"github.com/capamir/telegram-bot-go/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// StartHandler handles the /start command
func (b *Bot) StartHandler(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	welcomeText := `👋 Welcome to AI Bot!

I'm an intelligent assistant powered by Google Gemini AI, built with Go.

Commands:
/start - Show this message
/help - Get detailed help

💬 Just send me any message and I'll respond intelligently!

Examples:
• Ask me questions
• Request translations
• Get explanations
• Have a conversation`

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
Simply send me any text message and I'll respond using AI!

Available Commands:
/start - Welcome message
/help - This help message

What I Can Do:
• Answer questions on any topic
• Translate text between languages
• Explain complex concepts
• Write creative content
• Provide recommendations
• And much more!

Examples:
"What is quantum computing?"
"Translate 'hello' to Persian"
"Explain recursion simply"
"Tell me a joke"

Powered by Google Gemini 2.5 Flash 🤖✨`

	_, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpText,
		// ← Remove ParseMode line
	})
	if err != nil {
		log.Printf("Error sending help message: %v", err)
	}
}

// defaultHandler handles all non-command messages - now AI-powered!
func (b *Bot) defaultHandler(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	// Only process text messages
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	// Convert Telegram message to domain.Message
	domainMsg := &domain.Message{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
		// Tone, Thinking, Command will use default values (empty strings)
	}

	// Call usecase to handle message
	response, err := b.usecase.HandleMessage(ctx, domainMsg)
	if err != nil {
		log.Printf("Error handling message: %v", err)
		
		// Send error message to user
		_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Sorry, I encountered an error processing your message. Please try again.",
		})
		return
	}

	// Send AI response back to user
	_, err = tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   response.Text,
	})
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}
