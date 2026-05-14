# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Sudoku Daily API

Go 1.25 REST API using Fiber framework with PostgreSQL/Bun ORM.

## Project Structure

- `cmd/api/main.go` - API server entry point
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

# Run load tests (Vegeta + pghero)
make test-loads

# View load test reports
make load-report

# Clean load test reports
make load-clean

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

## Load Testing

Uses **Vegeta CLI** for load testing high-traffic endpoints. Activated via Docker Compose profiles.

### Endpoints Tested
| Endpoint | Method | Rate | Justification |
|----------|--------|------|---------------|
| `/api/sudoku` | GET | 1000 req/s | Highest volume - daily user access |
| `/api/sudoku/submit` | POST | 500 req/s | Authenticated submissions |
| `/api/sudoku/submit/guest` | POST | 500 req/s | Guest submissions |
| `/api/leaderboard` | GET | 300 req/s | Heavy queries with JOINs |
| `/api/auth/login` | POST | 200 req/s | Daily user logins |

### Running Load Tests
```bash
# Start load tests (uses docker-compose.load.yaml)
make test-loads

# View generated reports
make load-report

# Clean reports
make load-clean
```

Uses separate `docker-compose.load.yaml` to avoid interfering with the main development environment.

### Profiling (pprof)
- Activated when `LOG_LEVEL=debug` or `DEBUG=true`
- Runs on port `:6060` (internal Docker network only)
- Access from within Docker: `http://api:6060/debug/pprof/`
- Generate CPU profile: `go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30`

### Database Analysis
- **pghero**: Available at `http://localhost:8082` when using `make test-loads`
- **pg_stat_statements**: Enabled automatically via init script in `scripts/db/init.sql`
- Analyze slow queries in pghero dashboard during load tests
- Load test DB runs on port `5334` to avoid conflicts

### Load Test Reports
Reports are saved to `load-reports/` directory:
- Text reports (*.txt): Latency percentiles, throughput, errors
- HTML plots (*.html): Visual graphs of performance metrics

### Test Structure
```
scripts/
  load-tests/
    docker-compose.yaml     # Isolated load test environment
    run-load-tests.sh       # Main orchestrator
    setup-auth.sh           # Generates JWT tokens for auth endpoints
    targets/                # Vegeta target files
      get-sudoku.txt
      submit-guest.txt
      submit-auth.txt
      leaderboard.txt
      login.txt
  db/
    init.sql                # Enables pg_stat_statements
```

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