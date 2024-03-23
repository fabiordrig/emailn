include .env


test:
	@echo "Running tests..."
	go test -v -cover ./...

coverage:
	@echo "Running coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

build:
	@echo "Building..."
	go build -o bin/$(APP_NAME) main.go


PHONY: test coverage build
