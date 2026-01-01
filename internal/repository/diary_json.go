package repository

import (
    "context"
    "encoding/json"
    "errors"
    "os"
    "sync"

    "github.com/capamir/telegram-bot-go/internal/domain"
)

type JSONDiaryRepository struct {
    filePath string
    mu       sync.Mutex
}

func NewJSONDiaryRepository(filePath string) *JSONDiaryRepository {
    return &JSONDiaryRepository{
        filePath: filePath,
    }
}

// internal struct for file format
type diaryFile struct {
    Entries []*domain.DiaryEntry `json:"entries"`
}

func (r *JSONDiaryRepository) Save(ctx context.Context, entry *domain.DiaryEntry) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    data, err := r.load()
    if err != nil {
        return err
    }

    // check duplicate ID
    for _, e := range data.Entries {
        if e.ID == entry.ID {
            return ErrAlreadyExists
        }
    }

    data.Entries = append(data.Entries, entry)
    return r.save(data)
}

func (r *JSONDiaryRepository) GetByUser(ctx context.Context, userID int64) ([]*domain.DiaryEntry, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    data, err := r.load()
    if err != nil {
        return nil, err
    }

    var result []*domain.DiaryEntry
    for _, e := range data.Entries {
        if e.UserID == userID {
            result = append(result, e)
        }
    }
    return result, nil
}

func (r *JSONDiaryRepository) GetByID(ctx context.Context, id string) (*domain.DiaryEntry, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    data, err := r.load()
    if err != nil {
        return nil, err
    }

    for _, e := range data.Entries {
        if e.ID == id {
            return e, nil
        }
    }
    return nil, ErrNotFound
}

// helper: load file (create empty if not exists)
func (r *JSONDiaryRepository) load() (*diaryFile, error) {
    // if file doesn't exist, return empty
    _, err := os.Stat(r.filePath)
    if errors.Is(err, os.ErrNotExist) {
        return &diaryFile{Entries: []*domain.DiaryEntry{}}, nil
    } else if err != nil {
        return nil, err
    }

    b, err := os.ReadFile(r.filePath)
    if err != nil {
        return nil, err
    }
    if len(b) == 0 {
        return &diaryFile{Entries: []*domain.DiaryEntry{}}, nil
    }

    var data diaryFile
    if err := json.Unmarshal(b, &data); err != nil {
        return nil, err
    }
    if data.Entries == nil {
        data.Entries = []*domain.DiaryEntry{}
    }
    return &data, nil
}

// helper: save file
func (r *JSONDiaryRepository) save(data *diaryFile) error {
    b, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }
    // write atomically: write to temp then rename
    tmpPath := r.filePath + ".tmp"
    if err := os.WriteFile(tmpPath, b, 0o644); err != nil {
        return err
    }
    return os.Rename(tmpPath, r.filePath)
}
