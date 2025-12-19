package domain

import "time"

// DiaryEntry represents a saved diary entry with metadata
type DiaryEntry struct {
    ID           string    // UUID for this entry
    UserID       int64     // Telegram user ID (who wrote it)
    Username     string    // Telegram username (for display)
    ChatID       int64     // Where it was written (private or group)
    MessageID    int        // Telegram message ID (link back)
    OriginalText string    // User's diary text
    Summary      string    // AI summary with #diary tag
    Tone         Tone      // Tone used for summarization
    Mood         string    // AI-detected mood
    Tags         []string  // AI-generated topic tags
    CreatedAt    time.Time // When entry was created
}
