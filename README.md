# sudoku-daily-api

## Dependencies

- **Fiber**: REST API framework
- **Bun**: ORM for database
- **Postgres**: Database
- **Golang-migrate**: Database migrations
- **Viper**: Environment variables

## Endpoints

### Auth

| Method | Endpoint           | Description             |
| ------ | ------------------ | ----------------------- |
| POST   | /api/auth/register | Register a new user     |
| POST   | /api/auth/login    | Login and get tokens    |
| POST   | /api/auth/refresh  | Refresh access token    |
| POST   | /api/auth/logout   | Logout and revoke token |
| GET    | /api/auth/resume   | Get user statistics     |

### Sudoku

| Method | Endpoint             | Description             |
| ------ | -------------------- | ----------------------- |
| GET    | /api/sudoku          | Get daily sudoku        |
| POST   | /api/sudoku/generate | Generate daily sudokus  |
| POST   | /api/sudoku/submit   | Submit solution         |
| GET    | /api/sudoku/me       | Get user's daily solves |

### Leaderboard

| Method | Endpoint         | Description     |
| ------ | ---------------- | --------------- |
| GET    | /api/leaderboard | Get leaderboard |

## Response Codes

### Success Codes

| Code | Description                      |
| ---- | -------------------------------- |
| 200  | OK - Successful request          |
| 201  | Created - Resource created       |
| 204  | No Content - Successful deletion |

### Error Codes

| Code                     | HTTP Status | Description                                |
| ------------------------ | ----------- | ------------------------------------------ |
| invalid_query_param      | 400         | Invalid query parameter                    |
| invalid_email            | 401         | Invalid email format                       |
| invalid_token            | 401         | Invalid or missing token                   |
| token_expired            | 401         | Token has expired                          |
| invalid_credentials      | 401         | Invalid username or password               |
| email_already_registered | 409         | Email already exists                       |
| refresh_token_expired    | 401         | Refresh token has expired                  |
| refresh_token_revoked    | 401         | Refresh token was revoked                  |
| invalid_body             | 400         | Request body is invalid                    |
| invalid_solution         | 400         | Submitted solution is incorrect            |
| invalid_leaderboard_type | 400         | Invalid leaderboard type                   |
| size_required            | 400         | Size is required for this leaderboard type |
| size_not_allowed         | 400         | Size not allowed for this leaderboard type |
| invalid_size             | 400         | Invalid board size                         |
| invalid_limit            | 400         | Invalid limit value                        |
| invalid_page             | 400         | Invalid page number                        |
| internal_server_error    | 500         | Internal server error                      |
| too_many_requests        | 429         | Rate limit exceeded                        |
| already_played           | 409         | User already played today                  |
| user_not_found           | 404         | User not found                             |
| sudoku_not_found         | 404         | Sudoku not found                           |
| refresh_token_not_found  | 404         | Refresh token not found                    |
| solution_not_found       | 404         | Solution not found                         |
| validation_error         | 400         | Validation failed                          |

## Migrations

Migrations run on startup via `cmd/migrate/main.go`. Called by GitHub Actions pipeline.

## Running

```bash
# Start API
go run cmd/api/main.go

# Run migrations
go run cmd/migrate/main.go
```
