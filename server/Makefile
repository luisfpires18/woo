.PHONY: run build test lint clean gen-config

# Default target
run:
	go run ./cmd/server

# Build the server binary
build:
	go build -o bin/woo-server ./cmd/server

# Run all tests
test:
	go test ./... -v -count=1

# Run linter
lint:
	golangci-lint run ./...

# Format code
fmt:
	gofmt -w .

# Tidy dependencies
tidy:
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f data/woo.db

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Regenerate client JSON config from Go source of truth
gen-config:
	go run ./cmd/genconfig
