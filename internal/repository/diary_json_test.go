package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/capamir/telegram-bot-go/internal/domain"
)

func TestSaveAndGetByID(t *testing.T) {
	repo := NewJSONDiaryRepository(filepath.Join(t.TempDir(), "diary.json"))
	now := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)

	entry := &domain.DiaryEntry{
		ID:           "id-1",
		UserID:       42,
		Username:     "alice",
		ChatID:       100,
		MessageID:    200,
		OriginalText: "hello",
		Summary:      "summary",
		Mood:         "reflective",
		Tags:         []string{"tag1", "tag2"},
		CreatedAt:    now,
	}

	if err := repo.Save(context.Background(), entry); err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	got, err := repo.GetByID(context.Background(), "id-1")
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}

	if got.ID != entry.ID ||
		got.UserID != entry.UserID ||
		got.Username != entry.Username ||
		got.ChatID != entry.ChatID ||
		got.MessageID != entry.MessageID ||
		got.OriginalText != entry.OriginalText ||
		got.Summary != entry.Summary ||
		got.Mood != entry.Mood ||
		!got.CreatedAt.Equal(entry.CreatedAt) ||
		len(got.Tags) != len(entry.Tags) {
		t.Fatalf("stored entry mismatch: %#v", got)
	}
	for i := range entry.Tags {
		if got.Tags[i] != entry.Tags[i] {
			t.Fatalf("stored tags mismatch: %#v", got.Tags)
		}
	}

	if err := repo.Save(context.Background(), entry); !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestGetByUserFiltersResults(t *testing.T) {
	repo := NewJSONDiaryRepository(filepath.Join(t.TempDir(), "diary.json"))
	ctx := context.Background()

	entries := []*domain.DiaryEntry{
		{ID: "a1", UserID: 1, CreatedAt: time.Unix(1, 0)},
		{ID: "a2", UserID: 1, CreatedAt: time.Unix(2, 0)},
		{ID: "b1", UserID: 2, CreatedAt: time.Unix(3, 0)},
	}

	for _, entry := range entries {
		if err := repo.Save(ctx, entry); err != nil {
			t.Fatalf("save error: %v", err)
		}
	}

	got, err := repo.GetByUser(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestGetByUserEmpty(t *testing.T) {
	repo := NewJSONDiaryRepository(filepath.Join(t.TempDir(), "diary.json"))

	got, err := repo.GetByUser(context.Background(), 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestInvalidJSONReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "diary.json")
	if err := os.WriteFile(path, []byte("{invalid json"), 0o644); err != nil {
		t.Fatalf("write error: %v", err)
	}

	repo := NewJSONDiaryRepository(path)
	_, err := repo.GetByID(context.Background(), "id-1")
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
}

func TestConcurrentSaves(t *testing.T) {
	repo := NewJSONDiaryRepository(filepath.Join(t.TempDir(), "diary.json"))
	ctx := context.Background()

	const total = 10
	errs := make(chan error, total)
	var wg sync.WaitGroup

	for i := 0; i < total; i++ {
		wg.Add(1)
		entry := &domain.DiaryEntry{
			ID:        fmt.Sprintf("id-%d", i),
			UserID:    77,
			CreatedAt: time.Unix(int64(i), 0),
		}

		go func(e *domain.DiaryEntry) {
			defer wg.Done()
			errs <- repo.Save(ctx, e)
		}(entry)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("unexpected save error: %v", err)
		}
	}

	got, err := repo.GetByUser(ctx, 77)
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}
	if len(got) != total {
		t.Fatalf("expected %d entries, got %d", total, len(got))
	}
}
