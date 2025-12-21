// internal/utils/text.go
package utils

import (
    "fmt"
    "strings"
)

// FormatResponseWithQuestion prepends the user's question to the AI response
// and escapes HTML special characters
func FormatResponseWithQuestion(question, aiResponse string) string {
    // Escape HTML characters
    question = EscapeHTML(question)
    
    return fmt.Sprintf("❓ Question:\n%s\n\n📝 Answer:\n%s",
        question,
        aiResponse)
}

// EscapeHTML escapes HTML special characters for Telegram
func EscapeHTML(text string) string {
    replacer := strings.NewReplacer(
        "&", "&amp;",
        "<", "&lt;",
        ">", "&gt;",
    )
    return replacer.Replace(text)
}

// SplitMessage splits text into chunks without breaking words
// Useful for messages exceeding Telegram's 4096 character limit
func SplitMessage(text string, maxLength int) []string {
    if len(text) <= maxLength {
        return []string{text}
    }
    
    var parts []string
    remaining := text
    
    for len(remaining) > maxLength {
        splitAt := maxLength
        
        // Try to split at newline first
        if lastNewline := strings.LastIndex(remaining[:maxLength], "\n"); lastNewline > 0 {
            splitAt = lastNewline
        } else if lastSpace := strings.LastIndex(remaining[:maxLength], " "); lastSpace > 0 {
            // Fallback to space
            splitAt = lastSpace
        }
        
        parts = append(parts, strings.TrimSpace(remaining[:splitAt]))
        remaining = strings.TrimSpace(remaining[splitAt:])
    }
    
    if len(remaining) > 0 {
        parts = append(parts, remaining)
    }
    
    return parts
}
