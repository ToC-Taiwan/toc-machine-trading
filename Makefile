include .env
export

.PHONY: help

help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

swag-v1: ### swag init
	@echo 'package main' > ./tradebot.go && \
	swag init -g internal/controller/http/v1/router.go && \
	rm -rf ./tradebot.go
.PHONY: swag-v1

run: swag-v1 ### swag run
	@go mod tidy && go mod download && \
	CGO_ENABLED=0 go build -o toc-machine-trading ./cmd/app && ./toc-machine-trading
.PHONY: run

lint: ### check by golangci linter
	@golangci-lint run
.PHONY: lint

test: ### run test
	@go test -v -cover -race ./internal/...
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
