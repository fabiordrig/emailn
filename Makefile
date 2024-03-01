include .env.local


test:
	@echo "Running tests..."
	go test -v -cover ./...


