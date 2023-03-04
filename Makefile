include .env
export

BIN_NAME = toc-machine-trading

run: swag ### swag run
	@rm ./$(BIN_NAME)
	@go mod tidy
	@go mod download
	@go generate ./...
	@echo "Copying default config..."
	@cp ./configs/default.config.yml ./configs/config.yml
	@go build -o $(BIN_NAME) ./cmd/app
	@echo "Running $(BIN_NAME)..."
	@./$(BIN_NAME)
.PHONY: run

swag: ### swag
	@./scripts/generate_swagger.sh
.PHONY: swag

build: ### build
	@rm ./$(BIN_NAME)
	@go mod tidy
	@go mod download
	@go build -o $(BIN_NAME) ./cmd/app
.PHONY: build

go-mod-update: ### go-mod-update
	@./scripts/gomod_update.sh
.PHONY: go-mod-update

proto: ### proto
	@./scripts/compile_proto.sh
.PHONY: proto

update: go-mod-update proto swag ### update
.PHONY: update

lint: ### check by golangci linter
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run
.PHONY: lint

test: ### run test
	@go test ./... -v -coverprofile=coverage.txt -covermode=atomic
	@go tool cover -func coverage.txt
.PHONY: test

migrate-up-all: ### migration up to latest
	@migrate -path migrations -database '$(PG_URL)$(DB_NAME)?sslmode=disable' up
.PHONY: migrate-up-all

migrate-down-last: ### migration down one step
	@migrate -path migrations -database '$(PG_URL)$(DB_NAME)?sslmode=disable' down 1
.PHONY: migrate-down-last

migrate-create:  ### create new migration
	@migrate create -ext sql -dir migrations -tz "Asia/Taipei" -format "2006010215" 'migration'
.PHONY: migrate-create

help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help
