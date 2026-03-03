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
│   ├── 📂 migrate/				# migrations entrypoint
│   │   └── ▶️ main.go
├── 📂 migrations/				# For migration logic/files
│   ├── 📂 sql/					# sql up/down files
│   └── ▶️ migrate.go			# migration logic
├─ 📂 pkg/						# external logic like helpers
│   ├── 📂 config/				# configs initializing
│   └── 📂 database/			# database initialization
├── 📂 src/						# internal app files
│   ├── 📂 infrastructure/		# infrastructure
│   │   ├── 📁 persistence		# databases, queues
│   │   └── 📁 http				# Input ways http, grpc
│   ├── 📂 core/				# core domain structs and interfaces accessed by infrastructure
│   │   ├── 📁 models			# structs and their methods
│   │   └── 📁 ports			# interfaces
│   ├── 📁 application			# app logic called by main
│   └── 📁 services				# business logic
└── 📁 tests					# for tests (integration)
```

## Endpoints/Features

### Sudoku
- [x] GET /api/sudoku/daily?size=9
- [x] POST /api/sudoku/generate # Daily generation
- [ ] POST /api/sudoku/submit

### Auth
- [ ] POST /api/auth/register
- [ ] POST /api/auth/login
- [ ] GET /api/auth/me

## Migrations

The migrations run on the entrypoint cmd/migrate/main.go. It'll be called on Github Actions Pipe.
