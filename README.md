# Telegram AI Bot (Go)
**v0.2.0** - Clean Architecture MVP with AI Chat

A Telegram bot powered by Google Gemini AI, built with Clean Architecture in Go.

## Features

- 🤖 AI-powered chat using Google Gemini 2.5 Flash
- 🏗️ Clean Architecture (Domain, Usecase, Adapter layers)
- 🎭 Multiple tone support (Satirical, Serious, Friendly)
- ⚡ Context-aware responses with timeout handling

## Quick Start

### Prerequisites

- Go 1.21+
- Telegram Bot Token
- Google Gemini API Key

### Installation

1. Clone the repository
2. Copy `.env.example` to `.env`
3. Add your API keys to `.env`
4. Run: `go run cmd/bot/main.go`

## Commands

- `/start` - Welcome message
- `/help` - Usage instructions and examples
- `/ask <question>` - Ask AI a question
- `/settone <tone>` - Set response tone

## License

MIT
