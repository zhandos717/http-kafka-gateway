# Kafka Gateway Makefile

# Build variables
BINARY_NAME=kafka-gateway
MAIN_PATH=./cmd/kafkaGateway
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

kill:
	sudo kill -9 $(sudo lsof -t -i :8080)

# Build the project
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the project
run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH) && ./$(BUILD_DIR)/$(BINARY_NAME)

# Install dependencies
deps:
	$(GOMOD) tidy

# Test the project
test:
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./... && $(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

# Run go fmt
fmt:
	$(GOCMD) fmt ./...

# Run go vet
vet:
	$(GOCMD) vet ./...

# Build for different platforms
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

build-macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

# All builds
build-all: build-linux build-windows build-macos

# Help target
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  build           Build the project"
	@echo "  run             Build and run the project"
	@echo "  deps            Install dependencies"
	@echo "  test            Run tests"
	@echo "  test-coverage   Run tests with coverage report"
	@echo "  clean           Clean build artifacts"
	@echo "  fmt             Format code with go fmt"
	@echo "  vet             Run go vet"
	@echo "  build-linux     Build for Linux"
	@echo "  build-windows   Build for Windows"
	@echo "  build-macos     Build for macOS"
	@echo "  build-all       Build for all platforms"
	@echo "  help            Show this help message"

.PHONY: build run deps test test-coverage clean fmt vet build-linux build-windows build-macos build-all help