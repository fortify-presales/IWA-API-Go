# Makefile for a Go project
# Assumes entry point at project/cmd/server/server.go

# Binary output name
BINARY_NAME=insecure-go-api

# Default make command
all: build

# Build the binary
build:
	@echo "Building..."
	go build -o $(BINARY_NAME) ./cmd/server

# Run the server
run: build
	@echo "Running..."
	./$(BINARY_NAME)

# Test your application
test:
	@echo "Testing..."
	go test ./...

# Clean up binaries
clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)

# Generate mocks
mocks:
	@echo "Generating mocks..."
	mockgen -destination=internal/mocks/mock_repository.go  -package=mocks github.com/fortify-presales/insecure-go-api/internal/models Repository

# SAST Scan
sast-scan:
	@echo "Running SAST scan..."
#	gosec -exclude=G104 ./...
	sourceanalyzer -b insecure-go-api -clean
	sourceanalyzer -b insecure-go-api .
	sourceanalyzer -b insecure-go-api -scan 
	
# Help command to display available commands
help:
	@echo "Makefile commands:"
	@echo "all       - Build the application"
	@echo "build     - Build the binary"
	@echo "run       - Build and run the application"
	@echo "test      - Run tests"
	@echo "clean     - Remove binaries"
	@echo "mocks     - Generate mocks"
	@echo "sast-scan - Run SAST scan"
	@echo "help      - Display this help"

# Mark commands that don't correspond to files as .PHONY
.PHONY: all build run test clean help