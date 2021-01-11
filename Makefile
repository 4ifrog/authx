PROJECT_ROOT := $(shell pwd)
APP_BIN := ./bin
APP_NAME = authx
APP_SRC := ./cmd/$(APP_NAME)
TEST_SRC := ./pkg/...
IMAGE_NAME := cybersamx/$(APP_NAME)

# Target setup

.PHONY: all clean build help run docker lint

all: run

##@ run: Run application

run:
	@echo "Running $(APP_NAME)..."
	@cd $(APP_SRC); go run .

##@ build: Build binary

build:
	@-echo "Building $(APP_NAME)..."
	@-mkdir -p $(APP_BIN)
	go build -o $(APP_BIN) $(APP_SRC)
	@-cp $(APP_SRC)/config.yaml $(APP_BIN)

##@ docker-build: Build Docker image

docker:
	@echo "Building $(APP_NAME) docker image..."
	@docker \
		build \
		-t $(IMAGE_NAME) \
		.

##@ install: Install dependencies

install:
	@echo "Installing $(APP_NAME) dependencies..."
	go mod download

##@ lint: Run linter

lint:
	@echo "Linting $(APP_NAME)..."
	golangci-lint run -v

##@ format: Run gofmt

format:
	@echo "Formatting $(APP_NAME)..."
	gofmt -e -s -w .

##@ test: Run tests

test:
	@echo "Testing $(APP_NAME)..."
	go test $(TEST_SRC) -v

##@ clean: Clean output files and build cache

clean:
	@echo "Removing files and directories..."
	@-rm -rf $(APP_BIN)
	@-$(MAKE) go-clean

go-clean:
	@echo "Cleaning build cache..."
	go clean

##@ help: Help

.PHONY: help
all: help
help: Makefile
	@echo " Usage:\n  make <target>"
	@echo
	@sed -n 's/^##@//p' $< | column -t -s ':' | sed -e 's/[^ ]*/ &/2'
