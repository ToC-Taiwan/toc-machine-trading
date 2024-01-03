#!/bin/bash
set -e

echo "Generating swagger docs..."
go install github.com/swaggo/swag/cmd/swag@latest

echo 'package main' >./tmt.go
swag fmt -g internal/controller/http/router/router.go
swag init -q -g internal/controller/http/router/router.go
rm -rf ./tmt.go

echo "" >>./docs/swagger.json
git add ./docs

echo "Done"
