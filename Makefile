.PHONY: run build test tidy docker-up docker-down migrate lint

# Local development
run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

# Docker
docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api

# Database
migrate:
	@if [ -z "$(DB_URL)" ]; then \
		export $$(cat .env | xargs) && \
		psql "host=$$DB_HOST port=$$DB_PORT user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=$$DB_SSLMODE" \
			-f migrations/001_create_users.sql; \
	else \
		psql "$(DB_URL)" -f migrations/001_create_users.sql; \
	fi

# Linting (requires golangci-lint)
lint:
	golangci-lint run ./...
