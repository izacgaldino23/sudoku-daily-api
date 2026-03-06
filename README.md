# sudoku-daily-api

## Dependencies

**Fiber**: Rest api package manager;
**Bun**: To manage entities from the database;
**Postgres**: Database;
**Golang-migrate**: To manage database migrations;
**Viper**: To manage environment variables;

## Structure

```
📂 sudoku-daily-api/			# root
├── 📂 cmd/						# For main files
│   ├── 📂 api/					# api entrypoint
│   │   └── ▶️ main.go
│   └── 📂 migrate/				# migrations entrypoint
│       └── ▶️ main.go
├── 📂 migrations/				# For migration logic/files
│   ├── 📂 sql/					# sql up/down files
│   └── ▶️ migrate.go			# migration logic
├── 📂 pkg/						# external logic like helpers
│   ├── 📂 config/				# configs initializing
│   └── 📂 database/			# database initialization
├── 📂 src/						# internal app files
│   ├── 📂 domain/				# Domain layer
│   │   ├── 📂 entities/		# Domain entities
│   │   ├── 📂 vo/				# Value objects
│   │   ├── 📂 repository/		# Repository interfaces
│   │   ├── 📂 strategies/		# Generation algorithms
│   │   └── ▶️ services.go		# Domain service interfaces
│   ├── 📂 application/			# Use cases
│   │   └── 📂 usecase/
│   ├── 📂 infrastructure/		# Infrastructure layer
│   │   ├── 📂 http/			# HTTP handlers
│   │   └── 📂 persistence/		# Database persistence
│   └── 📂 services/			# Application services
└── 📁 tests					# for tests (integration)
```

## Endpoints/Features

### Sudoku
- [x] GET 	/api/sudoku/daily?size=9
- [x] POST 	/api/sudoku/generate # Daily generation
- [ ] POST 	/api/sudoku/submit

- [ ] Generate session token on GET /daily and validate on submit
- [ ] Validate solution on backend (compare with stored solution)

### Auth
- [x] POST 	/api/auth/register
- [x] POST 	/api/auth/login
- [x] POST 	/api/auth/refresh
- [ ] POST 	/api/auth/logout
- [ ] GET 	/api/auth/me

### Improvements

## Migrations

The migrations run on the entrypoint cmd/migrate/main.go. It'll be called on Github Actions Pipe.
