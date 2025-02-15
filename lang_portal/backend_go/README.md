# Lang Portal Backend

A Go-based backend server for the language learning portal application.

## Tech Stack

- Go with Gin web framework
- SQLite3 database
- Mage task runner

## Prerequisites

- Go 1.21 or higher
- SQLite3

## Project Structure

```
backend_go/
├── cmd/
│   ├── server/        # Main web server application
│   ├── init_db/       # Database initialization with CLI flags
│   ├── initdb/        # Simple database initialization
│   └── seed/          # Database seeding utility
├── internal/
│   ├── models/        # Data structures and database operations
│   ├── handlers/      # HTTP handlers organized by feature
│   ├── seeder/        # Seeding logic and data loading
│   └── services/      # Business logic
├── db/
│   ├── migrations/    # Database schema and migrations
│   └── seeds/         # Initial data for seeding
├── magefile.go        # Task runner configuration
├── go.mod            # Go module definition
└── go.sum            # Go module checksums
```

## Getting Started

1. Initialize the database:
```bash
mage initdb
```

2. Seed the database with initial data:
```bash
mage seed
```

3. Run the server:
```bash
mage run
```

## Available Mage Commands

- `mage build` - Build the application
- `mage run` - Run the server
- `mage initdb` - Initialize the database
- `mage seed` - Seed the database with initial data
- `mage test` - Run tests
- `mage clean` - Clean build artifacts
