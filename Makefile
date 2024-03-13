include .env.local


test:
	@echo "Running tests..."
	go test -v -cover ./...

coverage:
	@echo "Running coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

build:
	@echo "Building..."


PHONY: test coverage build
