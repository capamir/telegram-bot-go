package domain

import "time"

// DiaryEntry represents a saved diary entry with AI-generated insights
type DiaryEntry struct {
	// Unique identifiers
	ID        string    // UUID
	UserID    int64     // Telegram user ID (who wrote it)
	Username  string    // Telegram username (for display purposes)
	ChatID    int64     // Chat context (private or group)
	MessageID int       // Original Telegram message ID

	// Content
	OriginalText string   // User's raw diary text
	Summary      string   // AI-generated summary with #diary tag
	
	// AI-generated metadata
	Mood      string    // AI-detected mood (happy, sad, reflective, etc.)
	Tags      []string  // AI-generated topic tags
	
	// System metadata
	CreatedAt time.Time // Entry timestamp
}
