// internal/repository/diary.go
package repository

import (
    "context"
    "errors"
    
    "github.com/capamir/telegram-bot-go/internal/domain"
)

// Repository errors
var (
    ErrNotFound      = errors.New("diary entry not found")
    ErrAlreadyExists = errors.New("diary entry already exists")
)

// DiaryRepository defines storage operations for diary entries
type DiaryRepository interface {
    // Save creates a new diary entry
    // Returns ErrAlreadyExists if entry.ID already exists
    Save(ctx context.Context, entry *domain.DiaryEntry) error
    
    // GetByUser retrieves all diary entries for a user
    // Returns empty slice if user has no entries
    GetByUser(ctx context.Context, userID int64) ([]*domain.DiaryEntry, error)
    
    // GetByID retrieves a single diary entry by its ID
    // Returns ErrNotFound if entry doesn't exist
    GetByID(ctx context.Context, id string) (*domain.DiaryEntry, error)
}
