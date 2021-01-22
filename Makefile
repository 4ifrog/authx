PROJECT_ROOT := $(shell pwd)
APP_BIN := ./bin
APP_NAME = authx
APP_SRC := ./cmd/$(APP_NAME)
TEST_SRC := ./pkg/...
INT_TEST_SRC := ./cmd/integration-tests
IMAGE_NAME := cybersamx/$(APP_NAME)

# Target setup

.PHONY: all clean build help run docker lint test int-test int-test-full int-test-docker

all: run

##@ run: Run application

run:
	@echo "Running $(APP_NAME)..."
	@cd $(APP_SRC); go run .

##@ build: Build application

build:
	@-echo "Building $(APP_NAME)..."
	@-mkdir -p $(APP_BIN)
	CGO_ENABLED=0 go build -o $(APP_BIN) $(APP_SRC)
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
	CGO_ENABLED=0 go mod download

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
	@echo "Running unit tests..."
	CGO_ENABLED=0 go test $(TEST_SRC) -v -count=1

##@ int-test: Run integration tests on the local machine/container w/o depenencies
int-test:
	CGO_ENABLED=0 go test $(INT_TEST_SRC) -v -count=1

##@ int-test: Run integration tests on the local machine/container with dependencies as containers

int-test-full:
	@echo "Running integration tests with dependencies as containers..."
	@-docker-compose -f docker/docker-compose.test.yaml down --volumes
	@docker-compose -f docker/docker-compose.test.yaml up mongo &	# Run in the background
	# The application has a retry loop for connecting to the database.
	@$(MAKE) int-test
	@-docker-compose -f docker/docker-compose.test.yaml down --volumes

##@ int-test-docker: Run integration tests in a container (useful for CI)

int-test-docker:
	@echo "Running integration tests in a docker container..."
	@-docker-compose -f docker/docker-compose.test.yaml down --volumes
	@docker-compose -f docker/docker-compose.test.yaml up --build --abort-on-container-exit
	@-docker-compose -f docker/docker-compose.test.yaml down --volumes

##@ clean: Clean output files and build cache

clean:
	@echo "Removing files and directories..."
	@-rm -rf $(APP_BIN)
	@-$(MAKE) go-clean

go-clean:
	@echo "Cleaning build and test cache..."
	go clean -testcache

##@ help: Help

.PHONY: help
all: help
help: Makefile
	@echo " Usage:\n  make <target>"
	@echo
	@sed -n 's/^##@//p' $< | column -t -s ':' | sed -e 's/[^ ]*/ &/2'
