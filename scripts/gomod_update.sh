#!/bin/bash

rm go.mod
rm go.sum

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/golang/mock/mockgen@latest

go mod init github.com/toc-taiwan/toc-machine-trading
go mod tidy

git add go.mod go.sum
