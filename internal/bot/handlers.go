package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/capamir/telegram-bot-go/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// sendAIResponse handles loading state and sends AI response
func (b *Bot) sendAIResponse(
    ctx context.Context,
    tgBot *bot.Bot,
    chatID int64,
    questionText string,
    responseFunc func() (*domain.Response, error),
) {
    // 1. Show initial typing indicator
    _, _ = tgBot.SendChatAction(ctx, &bot.SendChatActionParams{
        ChatID: chatID,
        Action: models.ChatActionTyping,
    })
    
    // 2. Create channels
    type result struct {
        response *domain.Response
        err      error
    }
    responseChan := make(chan result, 1)
    stopTyping := make(chan struct{})
    
    // 3. Call AI in goroutine
    go func() {
        resp, err := responseFunc()
        responseChan <- result{response: resp, err: err}
    }()
    
    var loadingMsg *models.Message
    var response *domain.Response
    var err error
    
    // 4. Wait for response with timeout
    select {
    case res := <-responseChan:
        // Got response quickly (< 5s)
        response = res.response
        err = res.err
        
    case <-time.After(5 * time.Second):
        // Taking too long, send loading message
        loadingMsg, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
            ChatID: chatID,
            Text:   "⏳ Processing your request...",
        })
        
        // Keep showing typing indicator every 5s
        go func() {
            ticker := time.NewTicker(5 * time.Second)
            defer ticker.Stop()
            
            for {
                select {
                case <-ctx.Done():
                    return
                case <-stopTyping:
                    return
                case <-ticker.C:
                    _, _ = tgBot.SendChatAction(ctx, &bot.SendChatActionParams{
                        ChatID: chatID,
                        Action: models.ChatActionTyping,
                    })
                }
            }
        }()
        
        // Wait for actual response
        res := <-responseChan
        response = res.response
        err = res.err
        
        // Stop typing indicator
        close(stopTyping)
    }
    
    // 5. Handle errors
    if err != nil {
        log.Printf("Error getting AI response: %v", err)
        
        // Delete loading message if exists
        if loadingMsg != nil {
            _, _ = tgBot.DeleteMessage(ctx, &bot.DeleteMessageParams{
                ChatID:    chatID,
                MessageID: loadingMsg.ID,
            })
        }
        
        _, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
            ChatID: chatID,
            Text:   "Sorry, I encountered an error. Please try again.",
        })
        return
    }
    
    // 6. Format response
    formattedResponse := formatResponseWithQuestion(questionText, response.Text)
    
    // 7. Handle long responses (split if needed)
    if len(formattedResponse) > 4096 {
        // Response is too long, send in chunks
        b.sendLongMessage(ctx, tgBot, chatID, formattedResponse, loadingMsg)
        return
    }
    
    // 8. Send response as NEW message (don't edit)
    _, err = tgBot.SendMessage(ctx, &bot.SendMessageParams{
        ChatID:    chatID,
        Text:      formattedResponse,
        ParseMode: models.ParseModeHTML,
    })
    
    if err != nil {
        log.Printf("Error sending response: %v", err)
    }
    
    // 9. Delete loading message after sending response (optional)
    // Uncomment if you want to remove the "Processing..." message
    // if loadingMsg != nil {
    //     time.Sleep(500 * time.Millisecond)  // Brief delay so user sees transition
    //     _ = tgBot.DeleteMessage(ctx, &bot.DeleteMessageParams{
    //         ChatID:    chatID,
    //         MessageID: loadingMsg.ID,
    //     })
    // }
}

// sendLongMessage splits and sends messages that exceed Telegram's limit
func (b *Bot) sendLongMessage(
    ctx context.Context,
    tgBot *bot.Bot,
    chatID int64,
    text string,
    loadingMsg *models.Message,
) {
    const maxLength = 4000  // Leave buffer for formatting
    
    // Delete loading message first
    if loadingMsg != nil {
        _, _ = tgBot.DeleteMessage(ctx, &bot.DeleteMessageParams{
            ChatID:    chatID,
            MessageID: loadingMsg.ID,
        })
    }
    
    // Split message into chunks
    parts := splitMessage(text, maxLength)
    
    for i, part := range parts {
        // Add part indicator if multiple parts
        if len(parts) > 1 {
            part = fmt.Sprintf("📄 Part %d/%d\n\n%s", i+1, len(parts), part)
        }
        
        _, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
            ChatID:    chatID,
            Text:      part,
            ParseMode: models.ParseModeHTML,
        })
        
        if err != nil {
            log.Printf("Error sending message part %d: %v", i+1, err)
        }
        
        // Small delay between parts
        if i < len(parts)-1 {
            time.Sleep(300 * time.Millisecond)
        }
    }
}

// splitMessage splits text into chunks without breaking words
func splitMessage(text string, maxLength int) []string {
    if len(text) <= maxLength {
        return []string{text}
    }
    
    var parts []string
    remaining := text
    
    for len(remaining) > maxLength {
        // Find last newline before maxLength
        splitAt := maxLength
        lastNewline := strings.LastIndex(remaining[:maxLength], "\n")
        
        if lastNewline > 0 {
            splitAt = lastNewline
        } else {
            // Find last space before maxLength
            lastSpace := strings.LastIndex(remaining[:maxLength], " ")
            if lastSpace > 0 {
                splitAt = lastSpace
            }
        }
        
        parts = append(parts, strings.TrimSpace(remaining[:splitAt]))
        remaining = strings.TrimSpace(remaining[splitAt:])
    }
    
    if len(remaining) > 0 {
        parts = append(parts, remaining)
    }
    
    return parts
}


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

// formatResponseWithQuestion prepends the user's question to the AI response
func formatResponseWithQuestion(question, aiResponse string) string {
	// Escape only < > & for HTML
	question = strings.ReplaceAll(question, "&", "&amp;")
	question = strings.ReplaceAll(question, "<", "&lt;")
	question = strings.ReplaceAll(question, ">", "&gt;")

	return fmt.Sprintf("❓ <b>Question:</b>\n%s\n\n📝 <b>Answer:</b>\n%s",
		question,
		aiResponse)
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

	formattedResponse := formatResponseWithQuestion(questionText, response.Text)

	_, err = tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      formattedResponse,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Printf("Error sending movie-night response: %v", err)
	}
}
