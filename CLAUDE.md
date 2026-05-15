# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Sudoku Daily API

Go 1.25 REST API using Fiber framework with PostgreSQL/Bun ORM.

## Project Structure

- `cmd/api/main.go` - API server entry point (runs migrations on startup if enabled)
- `src/application/` - Fiber app setup, routes, bootstrap, and use cases
- `src/services/` - Business logic implementations (generator, password hasher, token, sudoku fetcher)
- `src/domain/` - Entities, repository interfaces, value objects, service interfaces, strategies
- `src/infrastructure/` - HTTP handlers, logging, database persistence
- `src/infrastructure/http/middlewares/metrics.go` - Prometheus metrics definitions + /metrics handler
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

# API server (runs migrations on startup if DATABASE_MIGRATIONS_ENABLED=true)
go run cmd/api/main.go
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
- Prometheus client_golang — /metrics endpoint + HTTP request metrics

## Environment Variables

- `LOG_LEVEL` — Log level (debug, info, warn, error, disabled). Default: info.
- `DEBUG` — Legacy toggle (true → debug level). Overridden by LOG_LEVEL.
- `DATABASE_MIGRATIONS_ENABLED` — Run migrations on API startup (true/false). Default: false.
- `DATABASE_MIGRATIONS_PATH` — Path to migration SQL files (e.g. `migrations/sql`).
- `AUTH_CRON_SECRET` — Secret for cron job authentication (via `Authorization: Bearer` or `X-Cron-Secret` header).
- `AUTH_OIDC_ENABLED` / `AUTH_OIDC_AUDIENCE` — OIDC authentication toggle and audience.

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
| POST   | /api/sudoku/generate/{size} | Generate daily sudokus (size: four, six, nine) |
| POST   | /api/sudoku/submit  | Submit solution (logged) |
| POST   | /api/sudoku/submit/guest | Submit solution (guest) |
| GET    | /api/sudoku/me     | Get user's daily solves |

### Leaderboard
| Method | Endpoint         | Description     |
| ------ | ---------------- | --------------- |
| GET    | /api/leaderboard | Get leaderboard |
| POST   | /api/leaderboard/reset | Reset daily strikes (OIDC/cron auth) |

### Cron
| Method | Endpoint               | Description                                     |
| ------ | ---------------------- | ----------------------------------------------- |
| POST   | /api/cron/generate/{size} | Generate daily sudoku for size (four, six, nine) |
| POST   | /api/leaderboard/reset | Reset strikes (also available via cron auth)    |

### Monitoring
| Method | Endpoint   | Description |
| ------ | ---------- | ----------- |
| GET    | /metrics   | Prometheus metrics (no auth) |
| GET    | /health    | Health check |

## Error Codes

Common error codes returned in responses: `invalid_token`, `token_expired`, `invalid_credentials`, `email_already_registered`, `invalid_solution`, `already_played`, `user_not_found`, `sudoku_not_found`, `internal_server_error`

---

**IMPORTANT**: Remember to update this file after modifying:
- Project structure
- Dependencies (go.mod changes)
- API endpoints or routes
- Running commands
- Environment variables

**DOCUMENTATION**: For every change/addition on handler, update the swagger docs
Run `make generate-docs`

## Workflows

- `.github/workflows/ci.yml` — Lint, build, unit + integration tests on PR/push to develop
- `.github/workflows/cd.yml` — Build + push Docker image to DockerHub on push to main
- `.github/workflows/cron.yml` — Daily cron at 00:00 UTC: generates all 3 sudoku sizes + resets leaderboard. Requires `API_URL` and `CRON_SECRET` secrets.`