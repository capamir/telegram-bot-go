package bot

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// StartHandler handles the /start command
func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	welcomeText := `👋 *Welcome to AI Bot!*

I'm an intelligent assistant powered by Google Gemini AI, built with Go.

*Commands:*
/start - Show this message
/help - Get detailed help

💬 *Just send me any message and I'll respond intelligently!*

Examples:
• Ask me questions
• Request translations
• Get explanations
• Have a conversation`

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      welcomeText,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("Error sending start message: %v", err)
	}
}

// HelpHandler handles the /help command
func HelpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	helpText := `📚 *Help & Information*

*How to Use:*
Simply send me any text message and I'll respond using AI!

*Available Commands:*
/start - Welcome message
/help - This help message

*What I Can Do:*
• Answer questions on any topic
• Translate text between languages
• Explain complex concepts
• Write creative content
• Provide recommendations
• And much more!

*Examples:*
"What is quantum computing?"
"Translate 'hello' to Persian"
"Explain recursion simply"
"Tell me a joke"

_Powered by Google Gemini 2.5 Flash_ 🤖✨`

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      helpText,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("Error sending help message: %v", err)
	}
}

// defaultHandler handles all non-command messages (echo for now)
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Only process text messages
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	// TODO: Phase 4 - Replace with AI-powered responses
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Echo: " + update.Message.Text + "\n\n_AI integration coming in Phase 4!_",
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("Error in default handler: %v", err)
	}
}
