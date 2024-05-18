BIN_NAME = toc-machine-trading

run: swag
	@go mod tidy
	@go mod download
	@go generate ./...
	@echo "Copying default config..."
	@cp ./configs/default.config.yml ./configs/config.yml
	@echo "Building $(BIN_NAME)..."
	@go build -o $(BIN_NAME) ./cmd/app
	@echo "Running $(BIN_NAME)..."
	@./$(BIN_NAME)

build:
	@go mod tidy
	@go mod download
	@go build -o $(BIN_NAME) ./cmd/app

swag:
	@./scripts/generate_swagger.sh

go-mod-update:
	@./scripts/gomod_update.sh

update: go-mod-update swag

lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run

test:
	@go test ./... -v -coverprofile=coverage.txt -covermode=atomic
	@go tool cover -func coverage.txt

migrate-up-all:
	@migrate -path migrations -database '$(PG_URL)$(DB_NAME)?sslmode=disable' up

migrate-down-last:
	@migrate -path migrations -database '$(PG_URL)$(DB_NAME)?sslmode=disable' down 1

migrate-create:
	@migrate create -ext sql -dir migrations -tz "Asia/Taipei" -format "2006010215" 'migration'

clean:
	@echo "Clean go cache..."
	@go clean -cache
	@go clean -modcache

help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
