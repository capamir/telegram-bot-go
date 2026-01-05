package usecase

// internal/usecase/chat.go

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/capamir/telegram-bot-go/internal/domain"
	"github.com/capamir/telegram-bot-go/internal/repository"
	"github.com/google/uuid"
)

// AIProvider defines what the use case needs from an AI model.
// This interface belongs to the USECASE layer.
type AIProvider interface {
    Generate(ctx context.Context, prompt string) (string, error)
}


// ChatUsecase contains application business logic.
type ChatUsecase struct {
    ai AIProvider
	diaryRepo repository.DiaryRepository
}


// NewChatUsecase injects dependencies into the use case.
func NewChatUsecase(ai AIProvider, repo repository.DiaryRepository) *ChatUsecase {
	return &ChatUsecase{
		ai:        ai,
		diaryRepo: repo,
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

// HandleDiary processes a diary entry, gets AI insights, and saves it
func (uc *ChatUsecase) HandleDiary(ctx context.Context, msg *domain.Message) (*domain.Response, error) {
	if msg.Text == "" {
		return nil, errors.New("diary text is empty")
	}

	// 1. Get structured AI response
	prompt := buildDiaryPrompt(msg.Text)
	aiRawResponse, err := uc.ai.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("ai error: %w", err)
	}

	// 2. Parse AI response (Expected to contain a JSON block)
	entryData, err := parseDiaryJSON(aiRawResponse)
	if err != nil {
		// Fallback: If JSON parsing fails, use raw response as summary
		entryData = &domain.DiaryEntry{
			Summary: aiRawResponse,
			Mood:    "Reflective",
			Tags:    []string{"diary"},
		}
	}

	// 3. Populate full domain entity
	entryData.ID = uuid.New().String()
	entryData.UserID = msg.ChatID // Using ChatID as UserID for simple private chats
	entryData.OriginalText = msg.Text
	entryData.CreatedAt = time.Now()

	// 4. Save to repository
	if uc.diaryRepo != nil {
		if err := uc.diaryRepo.Save(ctx, entryData); err != nil {
			return nil, fmt.Errorf("failed to save entry: %w", err)
		}
	}

	// 5. Format response for user
	responseText := fmt.Sprintf(
		"📔 <b>Diary Saved!</b>\n\n<b>Mood:</b> %s\n<b>Tags:</b> #%s\n\n<b>Summary:</b>\n%s",
		entryData.Mood,
		strings.Join(entryData.Tags, " #"),
		entryData.Summary,
	)

	return &domain.Response{Text: responseText}, nil
}

func parseDiaryJSON(raw string) (*domain.DiaryEntry, error) {
	// Simple extraction: find content between ```json and ```
	start := strings.Index(raw, "```json")
	if start == -1 {
		return nil, errors.New("no json block found")
	}
	content := raw[start+7:]
	end := strings.Index(content, "```")
	if end == -1 {
		return nil, errors.New("unclosed json block")
	}
	jsonStr := content[:end]

	var entry domain.DiaryEntry
	if err := json.Unmarshal([]byte(jsonStr), &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}
