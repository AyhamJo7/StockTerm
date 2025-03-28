# StockTerm Makefile

# Variables
BINARY_NAME=stockterm
BUILD_DIR=build
VERSION=1.0.0
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run

# Determine the operating system
ifeq ($(OS),Windows_NT)
	BINARY_SUFFIX=.exe
else
	BINARY_SUFFIX=
endif

.PHONY: all build clean test run tidy fmt lint help

all: clean build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)$(BINARY_SUFFIX) ./cmd/stockterm

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

run:
	@echo "Running $(BINARY_NAME)..."
	$(GORUN) ./cmd/stockterm

tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

help:
	@echo "Available commands:"
	@echo "  make build    - Build the binary"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make run      - Run the application"
	@echo "  make tidy     - Tidy dependencies"
	@echo "  make fmt      - Format code"
	@echo "  make lint     - Lint code"
	@echo "  make help     - Show this help message"
