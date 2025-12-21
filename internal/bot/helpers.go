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
    formattedResponse := utils.FormatResponseWithQuestion(questionText, response.Text)

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
    if loadingMsg != nil {
        time.Sleep(500 * time.Millisecond)  // Brief delay so user sees transition
        _, _ = tgBot.DeleteMessage(ctx, &bot.DeleteMessageParams{
            ChatID:    chatID,
            MessageID: loadingMsg.ID,
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
    const maxLength = 4000  // Leave buffer for formatting
    
    // Delete loading message first
    if loadingMsg != nil {
        _, _ = tgBot.DeleteMessage(ctx, &bot.DeleteMessageParams{
            ChatID:    chatID,
            MessageID: loadingMsg.ID,
        })
    }
    
    // Split message into chunks
    parts := utils.SplitMessage(text, maxLength)
    
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
