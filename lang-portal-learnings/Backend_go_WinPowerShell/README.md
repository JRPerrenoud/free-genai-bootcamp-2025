# Language Portal Backend

A Go-based API backend for the language learning portal. This service provides endpoints for managing vocabulary, study sessions, and learning progress.

## Tech Stack

- Go 1.21+
- Gin (Web Framework)
- SQLite3 (Database)
- Mage (Task Runner)

## Project Structure

```
lang_portal_go/
├── cmd/
│   └── server/         # Main application entry point
├── internal/
│   ├── models/         # Data structures and database operations
│   ├── handlers/       # HTTP handlers organized by feature
│   └── services/       # Business logic
├── db/
│   ├── migrations/     # Database migrations
│   └── seeds/          # Initial data population
├── magefile.go         # Task runner configuration
├── go.mod             # Go module file
└── words.db           # SQLite database
```

## Getting Started

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Build the project:
   ```bash
   mage build
   ```

3. Run the server:
   ```bash
   mage run
   ```

## Development

- Build: `mage build`
- Run: `mage run`
- Test: `mage test`
- Clean: `mage clean`
