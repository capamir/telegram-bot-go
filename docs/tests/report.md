# Test Report: Telegram Bot Go

## Coverage
- Usecase: input validation, default tone behavior, prompt construction, error propagation, movie-night prompt flow.
- Prompt logic: tone mapping, prompt construction, genre normalization, default genre behavior.
- Utils: HTML escaping, response formatting, message splitting rules and limits.
- Repository (JSON diary): save/get behavior, duplicate detection, invalid JSON errors, user filtering, concurrent saves.
- Config: required environment validation and default model selection.

## Why These Tests Matter
- Validate core business rules around tone defaults and prompt construction so AI output is deterministic and explainable.
- Ensure user input errors are rejected early (missing chat ID or message text).
- Confirm persistence rules (duplicate IDs, per-user filtering, and concurrency safety).
- Verify user-facing formatting and message splitting behavior to avoid Telegram length errors.
- Validate configuration requirements and safe defaults to prevent runtime startup failures.

## Risks Eliminated
- Silent acceptance of invalid chat inputs.
- Incorrect or inconsistent prompt templates for tone selection.
- Data loss or duplication in JSON diary storage during concurrent writes.
- Failure to handle corrupted persistence data.
- Misconfigured environment variables leading to undefined runtime behavior.

## Remaining Assumptions and Gaps
- Telegram API and Gemini API integrations are not covered; these require networked integration tests or contract tests with mocks.
- Bot handler behavior, long-message chunking, and typing indicator concurrency are not directly tested due to external Telegram client coupling.
- AI response text is not HTML-escaped; if it contains HTML, Telegram render behavior may differ and could pose formatting or security risks.
- Extremely long single-token words in `SplitMessage` can still exceed "no word breaks" expectations; behavior is not explicitly specified.

## How to Run
- `go test ./...`
- Optional race check: `go test -race ./...`
