# Project Variables
APP_NAME = api
BUILD_DIR = bin
MIGRATIONS_DIR = ./cmd/migrate
MIGRATION_BUILD = $(BUILD_DIR)/migrate
CMD_DIR = ./cmd/$(APP_NAME)
API_BIN = $(BUILD_DIR)/$(APP_NAME)
MAIN_FILE = $(CMD_DIR)/main.go
GO = go
FMT = fmt
GO_TEST_FLAGS = -v -race
GO_LINT_TOOL = golangci-lint
GO_LINT_FLAGS = run --fix
GO_BUILD_FLAGS = -ldflags "-s -w"
GO_FMT_ARGUMENTS = github.com/raganrrlaw/...
PKG_LIST := $(shell go list ./... | grep -v /vendor/)


## Default target: Build, test, and lint
all: lint test build

## Build the application binary
build:
	@echo "Building $(APP_NAME)..."
	@$(GO) build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Binary built at $(BUILD_DIR)/$(APP_NAME)."

## Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

## Run unit tests
test:
	@echo "Running tests..."
	@$(GO) test $(GO_TEST_FLAGS) $(PKG_LIST)
	@echo "Tests completed."

## Run code linter
lint:
	@echo "Running linter..."
	@$(GO_LINT_TOOL) $(GO_LINT_FLAGS)
	@echo "Linting completed."

## Clean up binaries and other generated files
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned."

## Update Go dependencies
deps:
	@echo "Tidying up and updating dependencies..."
	@$(GO) mod tidy
	@$(GO) mod vendor
	@echo "Dependencies updated."

## Compile the migration binary
migrate-build:
	@echo "Building migration binary..."
	@$(GO) build -o $(MIGRATION_BUILD) $(MIGRATIONS_DIR)
	@echo "Migration binary built at $(MIGRATION_BUILD)."

## Goose up
migrate-up: migrate-build
	@echo "Running migration..."
	@$(MIGRATION_BUILD) db/migrations up
	@echo "Migration completed."

## Goose up-to
migrate-up-to: migrate-build
	@echo "Running migration..."
	@$(MIGRATION_BUILD) db/migrations up-to 000003
	@echo "Migration completed."


## Goose down
migrate-down: migrate-build
	@echo "Running migration..."
	@$(MIGRATION_BUILD) db/migrations down
	@echo "Migration completed."

## Goose status
migrate-status: migrate
	@echo "Running migration status..."
	@$(MIGRATION_BUILD) db/migrations status

## Format the code
fmt:
	@echo "Formatting code..."
	@$(GO) $(FMT) $(GO_FMT_ARGUMENTS)
	@echo "Code formatted."
