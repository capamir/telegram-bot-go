package usecase

// internal/usecase/chat.go

import (
	"context"
	"errors"

	"github.com/capamir/telegram-bot-go/internal/domain"
)

// AIProvider defines what the use case needs from an AI model.
// This interface belongs to the USECASE layer.
type AIProvider interface {
    Generate(ctx context.Context, prompt string) (string, error)
}


// ChatUsecase contains application business logic.
type ChatUsecase struct {
    ai AIProvider
}


// NewChatUsecase injects dependencies into the use case.
func NewChatUsecase(ai AIProvider) *ChatUsecase {
    return &ChatUsecase{
        ai: ai,
    }
}


// HandleMessage orchestrates the chat flow.
func (uc *ChatUsecase) HandleMessage(
    ctx context.Context,
    msg *domain.Message,
) (*domain.Response, error) {

    // 1. Validate required fields
    if msg.ChatID == 0 {
        return nil, errors.New("chatID is required")
    }
    if msg.Text == "" {
        return nil, errors.New("text is required")
    }

    // 2. Apply defaults
    if msg.Tone == "" {
        msg.Tone = domain.ToneFriendly
    }

    // 3. Build prompt
    prompt := buildPrompt(msg)

    // 4. Call AI provider
    aiResponse, err := uc.ai.Generate(ctx, prompt)
    if err != nil {
        return nil, err
    }

    // 5. Create and return domain response
    return &domain.Response{
        Text: aiResponse,
    }, nil
}


// buildPrompt constructs the AI prompt based on domain rules.
func buildPrompt(msg *domain.Message) string {
    instruction := getToneInstruction(msg.Tone)
    
    // If instruction exists, format as system prompt + user question
    if instruction != "" {
        return instruction + "\n\nUser question: " + msg.Text
    }
    
    // Fallback to direct question
    return msg.Text
}

// HandleMovieNight builds a movie prompt and gets recommendations from AI.
func (uc *ChatUsecase) HandleMovieNight(ctx context.Context, rawGenres string, chatID int64) (*domain.Response, error) {
    // Build the movie-specific prompt
    prompt := buildMoviePrompt(rawGenres)

    // Reuse the same flow as HandleMessage but without tone logic
    aiResponse, err := uc.ai.Generate(ctx, prompt)
    if err != nil {
        return nil, err
    }

    return &domain.Response{
        Text: aiResponse,
    }, nil
}
