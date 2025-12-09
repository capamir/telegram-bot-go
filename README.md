# Telegram Bot with Google Gemini AI 🤖

An intelligent Telegram bot powered by Google's Gemini AI, built with Go.

## ✨ Features

- 🤖 AI-powered responses using Google Gemini 2.5 Flash
- 💬 Natural conversation handling
- ⚡ Fast and efficient (written in Go)
- 🔒 Secure configuration management
- 📦 Clean, modular architecture

## 📋 Prerequisites

- Go 1.21 or higher
- Telegram Bot Token (get from [@BotFather](https://t.me/botfather))
- Google Gemini API Key (get from [AI Studio](https://aistudio.google.com/apikey))

## 🚀 Installation

1. **Clone the repository:**
git clone https://github.com/capamir/telegram-bot-go.git
cd telegram-bot-go

2. **Install dependencies:**
go mod tidy

3. **Configure environment:**
cp .env.sample .env

Edit .env and add your tokens

4. **Run the bot:**
go run cmd/bot/main.go

## 🧪 Testing

Run integration tests:
go test ./test -v

Run with coverage:
go test ./test -v -cover

## 📁 Project Structure

telegram-bot-go/
├── cmd/bot/ # Application entry point
├── internal/ # Private application code
│ ├── ai/ # AI integration (Gemini)
│ ├── bot/ # Telegram bot logic
│ └── config/ # Configuration management
└── test/ # Integration tests

## 💬 Bot Commands

- `/start` - Start the bot and see welcome message
- `/help` - Get help and usage information

## 🛠️ Development Roadmap

- [x] Phase 1: Environment setup
- [x] Phase 2: Basic Telegram bot
- [x] Phase 3: Gemini AI integration (standalone)
- [ ] Phase 4: Full AI-powered bot
- [ ] Phase 5: Advanced features (context, rate limiting)

## 📝 Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `TELEGRAM_BOT_TOKEN` | Telegram bot token from BotFather | Yes |
| `GEMINI_API_KEY` | Google Gemini API key | Yes |
| `GEMINI_MODEL` | Gemini model name | No (default: gemini-2.5-flash) |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License

## 👤 Author

**capamir**

Built with ❤️ using Go and Google Gemini AI
