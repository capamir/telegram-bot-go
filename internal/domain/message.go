package domain

// Enums
type Tone string
type ThinkingMode string
type CommandType string

const (
	ToneSatirical    Tone = "satirical"
	ToneSerious      Tone = "serious"
	ToneFriendly     Tone = "friendly"
	ToneProfessional Tone = "professional"
)

const (
	ThinkingStandard ThinkingMode = "standard"
	ThinkingExtended ThinkingMode = "extended"
)

const (
	CommandText         CommandType = "text"
	CommandFileAnalysis CommandType = "file_analysis"
	CommandSTT          CommandType = "voice_to_text"
	CommandTTS          CommandType = "text_to_voice"
)

// Entities
type Message struct {
	ChatID    int64
	Text      string
	Tone      Tone
	Thinking  ThinkingMode
	Command   CommandType
	MessageID int64
}

type Response struct {
	Text string
}
