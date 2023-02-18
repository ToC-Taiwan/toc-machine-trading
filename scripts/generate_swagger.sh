#!/bin/bash

echo "Generating swagger docs..."
echo 'package main' > ./tradebot.go
swag init -g internal/controller/http/v1/router.go
rm -rf ./tradebot.go
echo "" >> ./docs/swagger.json
git add ./docs
