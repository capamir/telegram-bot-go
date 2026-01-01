package usecase

import (
	"fmt"
	"testing"

	"github.com/capamir/telegram-bot-go/internal/domain"
)

func TestGetToneInstructionMapping(t *testing.T) {
	tests := []struct {
		name     string
		tone     domain.Tone
		expected string
	}{
		{name: "satirical", tone: domain.ToneSatirical, expected: promptSatirical},
		{name: "serious", tone: domain.ToneSerious, expected: promptSerious},
		{name: "friendly", tone: domain.ToneFriendly, expected: promptFriendly},
		{name: "professional", tone: domain.ToneProfessional, expected: promptProfessional},
		{name: "unknown defaults friendly", tone: domain.Tone("unknown"), expected: promptFriendly},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getToneInstruction(tt.tone)
			if got != tt.expected {
				t.Fatalf("unexpected instruction for %q", tt.name)
			}
		})
	}
}

func TestBuildPromptUsesInstruction(t *testing.T) {
	msg := &domain.Message{
		Text: "Explain Go interfaces",
		Tone: domain.ToneSerious,
	}

	expected := promptSerious + "\n\nUser question: " + msg.Text
	if got := buildPrompt(msg); got != expected {
		t.Fatalf("unexpected prompt:\nexpected: %q\nactual:   %q", expected, got)
	}
}

func TestBuildMoviePromptDefaultGenres(t *testing.T) {
	const defaultGenres = "Action Thriller, Crime/Heist, Mystery/Detective (action-oriented), Epic Fantasy, Thriller/Suspense, Adventure"
	expected := fmt.Sprintf("%s\n\nUser requested genres: %s", promptMovieNight, defaultGenres)

	if got := buildMoviePrompt(""); got != expected {
		t.Fatalf("unexpected default movie prompt:\nexpected: %q\nactual:   %q", expected, got)
	}
}

func TestBuildMoviePromptCleansWhitespace(t *testing.T) {
	expected := fmt.Sprintf("%s\n\nUser requested genres: %s", promptMovieNight, "action thriller sci-fi")

	if got := buildMoviePrompt("  action   thriller \n sci-fi "); got != expected {
		t.Fatalf("unexpected cleaned movie prompt:\nexpected: %q\nactual:   %q", expected, got)
	}
}
