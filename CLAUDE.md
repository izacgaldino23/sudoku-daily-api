# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Sudoku Daily API

Go 1.25 REST API using Fiber framework with PostgreSQL/Bun ORM.

## Project Structure

- `cmd/api/main.go` - API server entry point
- `cmd/migrate/main.go` - Database migrations entry point
- `src/application/` - Fiber app setup, routes, bootstrap, and use cases
- `src/services/` - Business logic implementations (generator, password hasher, token, sudoku fetcher)
- `src/domain/` - Entities, repository interfaces, value objects, service interfaces, strategies
- `src/infrastructure/` - HTTP handlers, logging, database persistence
- `pkg/` - Shared packages (config, database, errors, validator)
- `migrations/sql/*.sql` - Database migrations

## Architecture

The project follows a layered architecture:
- **Application Layer** (`src/application/`) - Route registration, middleware, use cases
- **Domain Layer** (`src/domain/`) - Business entities, repository interfaces, domain strategies
- **Infrastructure Layer** (`src/infrastructure/`) - HTTP handlers, database implementations
- **Services Layer** (`src/services/`) - Business logic (auth, sudoku generation, tokens)
- **Package Layer** (`pkg/`) - Cross-cutting concerns (config, errors, validation)

## Running

```bash
# API server
go run cmd/api/main.go

# Run tests (requires Docker)
go test ./tests/...

# Migrations
go run cmd/migrate/main.go
```

## Docker

- `docker-compose.yaml` - Local PostgreSQL
- Run `docker compose up -d` before starting API

## Make Commands

```bash
# Create new migration
make new-migration name=<migration_name>

# Run migrations
make run-migrations

# Run integration tests
make test-integration

# Generate API docs (Swagger)
make generate-docs
```

## Dependencies

- Fiber (HTTP), Bun (ORM), golang-jwt/jwt/v5, pg (postgres driver), Viper (env), golang-migrate, zerolog

## API Endpoints

### Auth
| Method | Endpoint           | Description             |
| ------ | ------------------ | ----------------------- |
| POST   | /api/auth/register | Register a new user     |
| POST   | /api/auth/login    | Login and get tokens    |
| POST   | /api/auth/refresh  | Refresh access token   |
| POST   | /api/auth/logout   | Logout and revoke token |
| GET    | /api/auth/resume   | Get user statistics     |

### Sudoku
| Method | Endpoint               | Description               |
| ------ | ----------------------| ----------------------- |
| GET    | /api/sudoku          | Get daily sudoku        |
| POST   | /api/sudoku/generate | Generate daily sudokus  |
| POST   | /api/sudoku/submit  | Submit solution (logged) |
| POST   | /api/sudoku/submit/guest | Submit solution (guest) |
| GET    | /api/sudoku/me     | Get user's daily solves |

### Leaderboard
| Method | Endpoint         | Description     |
| ------ | ---------------- | --------------- |
| GET    | /api/leaderboard | Get leaderboard |

## Error Codes

Common error codes returned in responses: `invalid_token`, `token_expired`, `invalid_credentials`, `email_already_registered`, `invalid_solution`, `already_played`, `user_not_found`, `sudoku_not_found`, `internal_server_error`

---

**IMPORTANT**: Remember to update this file after modifying:
- Project structure
- Dependencies (go.mod changes)
- API endpoints or routes
- Running commands
- Environment variables