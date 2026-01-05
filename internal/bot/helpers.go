// internal/bot/helpers.go
package bot

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/capamir/telegram-bot-go/internal/domain"
	"github.com/capamir/telegram-bot-go/internal/utils"
    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
)

// sendAIResponse handles loading state, background typing indicators, and sending AI responses.
// It uses a child context to ensure background workers are cleaned up immediately upon completion.
func (b *Bot) sendAIResponse(
	ctx context.Context,
	tgBot *bot.Bot,
	chatID int64,
	questionText string,
	responseFunc func() (*domain.Response, error),
) {
	// 1. Create a child context for background workers (typing indicator)
	// This ensures that when we finish (either success or error), the background goroutines stop.
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 2. Initial "Typing..." action
	_, _ = tgBot.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: chatID,
		Action: models.ChatActionTyping,
	})

	type result struct {
		response *domain.Response
		err      error
	}
	responseChan := make(chan result, 1)

	// 3. Start AI processing in a goroutine
	go func() {
		resp, err := responseFunc()
		responseChan <- result{response: resp, err: err}
	}()

	var loadingMsg *models.Message
	var response *domain.Response
	var err error

	// 4. Wait for response with a "long-running" fallback
	select {
	case res := <-responseChan:
		// Fast response: Received before the 5s threshold
		response = res.response
		err = res.err

	case <-time.After(5 * time.Second):
		// Slow response: Send a "Processing" message and keep the typing indicator alive
		loadingMsg, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "⏳ Processing your request...",
		})

		// Continuous typing indicator in background using the child context
		go func(c context.Context) {
			ticker := time.NewTicker(4 * time.Second) // Telegram typing lasts ~5s
			defer ticker.Stop()
			for {
				select {
				case <-c.Done():
					return
				case <-ticker.C:
					_, _ = tgBot.SendChatAction(c, &bot.SendChatActionParams{
						ChatID: chatID,
						Action: models.ChatActionTyping,
					})
				}
			}
		}(childCtx)

		// Block until the actual response arrives
		res := <-responseChan
		response = res.response
		err = res.err
		
		// Signal background workers to stop immediately
		cancel()
	}

	// 5. Handle potential errors
	if err != nil {
		log.Printf("AI Error for ChatID %d: %v", chatID, err)
		b.cleanupLoadingMessage(ctx, tgBot, chatID, loadingMsg)
		
		_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "⚠️ Sorry, I encountered an issue while processing your request. Please try again later.",
		})
		return
	}

	// 6. Format and deliver the response
	formattedResponse := utils.FormatResponseWithQuestion(questionText, response.Text)

	if len(formattedResponse) > 4000 {
		b.sendLongMessage(ctx, tgBot, chatID, formattedResponse, loadingMsg)
		return
	}

	b.cleanupLoadingMessage(ctx, tgBot, chatID, loadingMsg)

	_, err = tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      formattedResponse,
		ParseMode: models.ParseModeHTML,
	})
	
	if err != nil {
		log.Printf("Error sending message to %d: %v", chatID, err)
	}
}

// cleanupLoadingMessage deletes the "Processing..." message if it exists
func (b *Bot) cleanupLoadingMessage(ctx context.Context, tgBot *bot.Bot, chatID int64, msg *models.Message) {
	if msg != nil {
		_, _ = tgBot.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    chatID,
			MessageID: msg.ID,
		})
	}
}

// sendLongMessage splits and sends messages that exceed Telegram's limit
func (b *Bot) sendLongMessage(
	ctx context.Context,
	tgBot *bot.Bot,
	chatID int64,
	text string,
	loadingMsg *models.Message,
) {
	const maxLength = 3800 // Buffer for formatting and parts indicator
	
	b.cleanupLoadingMessage(ctx, tgBot, chatID, loadingMsg)

	parts := utils.SplitMessage(text, maxLength)
	for i, part := range parts {
		header := ""
		if len(parts) > 1 {
			header = fmt.Sprintf("📄 <b>Part %d/%d</b>\n\n", i+1, len(parts))
		}

		_, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      header + part,
			ParseMode: models.ParseModeHTML,
		})

		if err != nil {
			log.Printf("Error sending part %d to %d: %v", i+1, chatID, err)
		}
		
		// Prevent flooding
		time.Sleep(200 * time.Millisecond)
	}
}
