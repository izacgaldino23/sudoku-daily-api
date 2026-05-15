MIGRATIONS_PATH = migrations/sql
DATABASE_HOST = localhost
DATABASE_PORT = 5432
DATABASE_USER = postgres
DATABASE_PASSWORD = postgres
DATABASE_NAME = sudoku_db

# run migrate command with a custom name from args
new-migration:
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

run-migrations:
	set ENV=local&& DATABASE_MIGRATIONS_ENABLED=true go run cmd/api/main.go

# run unit tests for pkg and src folder on root
test:
	go test ./pkg/... ./src/... -p=1

# run integration tests using build tags
test-integration:
	go test -tags=integration ./tests/... -p=1

clean-test:
	docker-compose -f tests/docker-compose.test.yaml down

generate-docs:
	swag init -g ./cmd/api/main.go

lint: format
	golangci-lint run

format:
	go fmt ./...

docker-build:
	docker compose -f docker-compose.yaml up -d --remove-orphans --build

.PHONY: new-migration run-migrations test-integration generate-docs lint format test test-loads docker-build