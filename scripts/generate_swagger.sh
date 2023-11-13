#!/bin/bash

echo "Generating swagger docs..."
echo 'package main' >./tradebot.go
swag fmt
swag init --generatedTime=true -q -g internal/controller/http/router/router.go
rm -rf ./tradebot.go
echo "" >>./docs/swagger.json
git add ./docs
