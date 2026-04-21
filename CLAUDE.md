# Sudoku Daily API

Go 1.25 REST API using Fiber framework with PostgreSQL/Bun ORM.

## Project Structure

- `cmd/api/main.go` - API server entry point
- `cmd/migrate/main.go` - Database migrations entry point
- `src/application/app.go` - Fiber app setup and routes
- `src/services/` - Business logic (auth, sudoku, leaderboard, tokens)
- `pkg/` - Shared packages (config, errors, database)
- `migrations/sql/*.sql` - Database migrations

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

## Dependencies

- Fiber (HTTP), Bun (ORM), golang-jwt/jwt/v5, pg (postgres driver), Viper (env)

---

**IMPORTANT**: Remember to update this file after modifying:
- Project structure
- Dependencies (go.mod changes)
- API endpoints or routes
- Running commands
- Environment variables