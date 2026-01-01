package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/capamir/telegram-bot-go/internal/domain"
)

type stubAI struct {
	prompt   string
	response string
	err      error
}

func (s *stubAI) Generate(ctx context.Context, prompt string) (string, error) {
	s.prompt = prompt
	return s.response, s.err
}

func TestHandleMessageRequiresChatID(t *testing.T) {
	uc := NewChatUsecase(&stubAI{})

	_, err := uc.HandleMessage(context.Background(), &domain.Message{Text: "hello"})
	if err == nil || err.Error() != "chatID is required" {
		t.Fatalf("expected chatID error, got %v", err)
	}
}

func TestHandleMessageRequiresText(t *testing.T) {
	uc := NewChatUsecase(&stubAI{})

	_, err := uc.HandleMessage(context.Background(), &domain.Message{ChatID: 10})
	if err == nil || err.Error() != "text is required" {
		t.Fatalf("expected text error, got %v", err)
	}
}

func TestHandleMessageDefaultsToneAndBuildsPrompt(t *testing.T) {
	stub := &stubAI{response: "ok"}
	uc := NewChatUsecase(stub)

	msg := &domain.Message{
		ChatID: 1,
		Text:   "What is Go?",
	}

	resp, err := uc.HandleMessage(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Tone != domain.ToneFriendly {
		t.Fatalf("expected default tone to be friendly, got %q", msg.Tone)
	}

	expectedPrompt := promptFriendly + "\n\nUser question: " + msg.Text
	if stub.prompt != expectedPrompt {
		t.Fatalf("unexpected prompt:\nexpected: %q\nactual:   %q", expectedPrompt, stub.prompt)
	}
	if resp == nil || resp.Text != "ok" {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestHandleMessagePropagatesAIError(t *testing.T) {
	stub := &stubAI{err: errors.New("boom")}
	uc := NewChatUsecase(stub)

	resp, err := uc.HandleMessage(context.Background(), &domain.Message{
		ChatID: 2,
		Text:   "test",
	})
	if err == nil || err.Error() != "boom" {
		t.Fatalf("expected error propagation, got %v", err)
	}
	if resp != nil {
		t.Fatalf("expected nil response on error, got %#v", resp)
	}
}

func TestHandleMovieNightUsesPrompt(t *testing.T) {
	stub := &stubAI{response: "list"}
	uc := NewChatUsecase(stub)

	rawGenres := " action   thriller \n mystery "
	expectedGenres := "action thriller mystery"
	expectedPrompt := fmt.Sprintf("%s\n\nUser requested genres: %s", promptMovieNight, expectedGenres)

	resp, err := uc.HandleMovieNight(context.Background(), rawGenres, 99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.prompt != expectedPrompt {
		t.Fatalf("unexpected movie prompt:\nexpected: %q\nactual:   %q", expectedPrompt, stub.prompt)
	}
	if resp == nil || resp.Text != "list" {
		t.Fatalf("unexpected response: %#v", resp)
	}
}
