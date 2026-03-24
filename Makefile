.PHONY: new-migration run-migrations test-integration

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
	set ENV=local&& go run cmd/migrate/main.go

test-integration:
	docker-compose -f tests/docker-compose.test.yaml up -d
	go test ./tests/integration/... -v
	docker-compose -f tests/docker-compose.test.yaml down

generate-docs:
	swag init -g ./cmd/api/main.go