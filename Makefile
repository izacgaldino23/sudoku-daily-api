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

# run unit tests for pkg and src folder on root
test:
	go test ./pkg/... ./src/... -p=1

# run integration tests on tests folder on root
test-integration:
	docker-compose -f tests/docker-compose.test.yaml up -d 
	go test ./tests/integration/... -p=1

clean-test:
	docker-compose -f tests/docker-compose.test.yaml down

test-loads:
	docker-compose -f scripts/load-tests/docker-compose.yaml up -d --remove-orphans db-load pghero
	docker-compose -f scripts/load-tests/docker-compose.yaml run migrate-load
	docker-compose -f scripts/load-tests/docker-compose.yaml up -d --remove-orphans api-load
	docker-compose -f scripts/load-tests/docker-compose.yaml run vegeta
	docker-compose -f scripts/load-tests/docker-compose.yaml down --remove-orphans

load-report:
	@echo "Load test reports:"
	@ls -la load-reports/ 2>/dev/null || echo "No reports found. Run 'make test-loads' first."

load-clean:
	rm -f load-reports/*

generate-docs:
	swag init -g ./cmd/api/main.go

lint: format
	golangci-lint run

format:
	go fmt ./...

.PHONY: new-migration run-migrations test-integration generate-docs lint format test test-loads