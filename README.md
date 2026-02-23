# sudoku-daily-api

## Dependencies

**Fiber**: Rest api package manager;
**Gorm**: To manage entities from the database;

## Structure

```
📂 sudoku-daily-api/			# root
├── 📂 app/					# internal app files
│   ├── 📂 adapters/			# adapters
│   │   ├── 📁 drivens			# databases, queues
│   │   └── 📁 drivers			# Input ways rest, grpc
│   ├── 📂 core/				# core domain structs and interfaces accessed by adapters
│   │   ├── 📁 models			# structs and their methods
│   │   └── 📁 ports			# interfaces
│   ├── 📁 application			# app logic called by main
│   └── 📁 services			# business logic
├── 📁 pkg					# external logic like helpers
├── 📁 tests				# for tests (integration)
└── ▶️ main.go
```

## Endpoints

### Sudoku
- [ ] GET /api/sudoku/daily?size=9
- [ ] POST /api/sudoku/submit

### Auth
- [ ] POST /api/auth/register
- [ ] POST /api/auth/login
- [ ] GET /api/auth/me
