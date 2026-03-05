.PHONY: run build test tidy lint swag-init \
        docker-up docker-down docker-logs \
        migration-up migration-down \
        build-image push-image \
        rename-module

CURRENT_DIR=$(shell pwd)
APP=$(shell basename ${CURRENT_DIR})
APP_CMD_DIR=${CURRENT_DIR}/cmd/api
TAG=latest
ENV_TAG=latest

# Local development
run:
	go run ./cmd/main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

test:
	go test ./...

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

swag-init:
	swag init -g api/api.go -o api/docs --parseDependency --parseInternal

# Module rename — usage: make rename-module NEW_MODULE=github.com/your-username/my-project
rename-module:
	@if [ -z "$(NEW_MODULE)" ]; then \
		echo "Error: NEW_MODULE is required."; \
		echo "Usage: make rename-module NEW_MODULE=github.com/your-username/my-project"; \
		exit 1; \
	fi
	@OLD_MODULE=$$(go mod edit -json | grep '"Module"' -A1 | grep '"Path"' | sed 's/.*"Path": "\(.*\)".*/\1/'); \
	echo "Renaming $$OLD_MODULE → $(NEW_MODULE)"; \
	find . -type f -name "*.go" \
		-not -path "./vendor/*" \
		| xargs sed -i '' "s|$$OLD_MODULE|$(NEW_MODULE)|g"; \
	sed -i '' "s|$$OLD_MODULE|$(NEW_MODULE)|g" go.mod; \
	go mod tidy; \
	echo "Done. Module is now $(NEW_MODULE)."

# Docker
docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api

# Database migrations (requires golang-migrate: https://github.com/golang-migrate/migrate)
migration-up:
	migrate -path ./migrations/postgres -database '${DB_URL}' up

migration-down:
	migrate -path ./migrations/postgres -database '${DB_URL}' down

# Docker image
build-image:
	docker build --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} .
	docker tag ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

push-image:
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}
