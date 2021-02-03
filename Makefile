# Project
PROJECT_ROOT := $(shell pwd)
PROJECT_BIN := ./bin
APP_NAME = authx
APP_SRC := ./cmd/$(APP_NAME)
TEST_SRC := ./pkg/...
SRC_STATIC_WEB := ./web_client
TARGET_STATIC_WEB := ./web_client/build

# Deployment
IMAGE_NAME := cybersamx/$(APP_NAME)

# Colorized print
BOLD := $(shell tput bold)
RED := $(shell tput setaf 1)
BLUE := $(shell tput setaf 4)
CYAN := $(shell tput setaf 6)
RESET := $(shell tput sgr0)

# Set up the default target all and set up phony targets.

.PHONY: all

all: run

##@ run: Run application

.PHONY: run

run: web-build
	@-echo "$(BOLD)$(BLUE)Running $(APP_NAME)...$(RESET)"
	@cd $(APP_SRC); go run .

##@ build: Build application

.PHONY: build

build: web-build
	@-echo "$(BOLD)$(BLUE)Building $(APP_NAME)...$(RESET)"
	@mkdir -p $(PROJECT_BIN)
	CGO_ENABLED=0 go build -o $(PROJECT_BIN) $(APP_SRC)
	@cp $(APP_SRC)/config.yaml $(PROJECT_BIN)
	@mkdir -p $(PROJECT_BIN)/$(TARGET_STATIC_WEB)
	@cp -rf $(SRC_STATIC_WEB)/build/* $(PROJECT_BIN)/$(TARGET_STATIC_WEB)

##@ web-build: Build the web application.

.PHONY: web-build

web-build:
	@-echo "$(BOLD)$(BLUE)Building web application...$(RESET)"
	@cd $(SRC_STATIC_WEB) && \
	npm install && \
	npm run build && \
	cd -

##@ web-build: Test the web application.

.PHONY: web-test

web-test:
	@-echo "$(BOLD)$(BLUE)Building web application...$(RESET)"
	@cd $(SRC_STATIC_WEB) && \
	npm test && \
	cd -

##@ docker-build: Build Docker image

.PHONY: docker

docker:
	@-echo "$(BOLD)$(BLUE)Building $(APP_NAME) docker image...$(RESET)"
	@docker \
		build \
		-t $(IMAGE_NAME) \
		.

##@ lint: Run linter

.PHONY: lint

lint:
	@-echo "$(BOLD)$(BLUE)Linting $(APP_NAME)...$(RESET)"
	golangci-lint run -v
	@cd $(SRC_STATIC_WEB) && \
	npm run lint && \
	cd -

##@ format: Run gofmt

.PHONY: format

format:
	@-echo "$(BOLD)$(BLUE)Formatting $(APP_NAME)...$(RESET)"
	gofmt -e -s -w .
	@cd $(SRC_STATIC_WEB) && \
	npm run lint:fix && \
	cd -

##@ test: Run tests

.PHONY: test

test: start-db-containers web-test
	@-echo "$(BOLD)$(CYAN)Running tests...$(RESET)"
	CGO_ENABLED=0 go test $(TEST_SRC) -v -count=1 -coverprofile cover.out
	go tool cover -func cover.out

##@ int-containers: Run tests and databases as containers within a netwwork context (useful for CI)

.PHONY: test-containers

test-containers: start-db-containers
	@-echo "$(BOLD)$(CYAN)Running tests and dependencies as docker containers...$(RESET)"
	@docker-compose -f docker/docker-compose.test.yaml up --build --abort-on-container-exit

##@ start-db-containers: Start database containers if they aren't running in the background

.PHONY: start-db-containers

start-db-containers: scripts/start-db-containers.sh
	@-echo "$(BOLD)$(BLUE)Starting database containers...$(RESET)"
	$(PROJECT_ROOT)/scripts/start-db-containers.sh

##@ end-db-containers: End database containers if they are running in the background

.PHONY: end-db-containers

end-db-containers:
	@-echo "$(BOLD)$(BLUE)Ending database containers...$(RESET)"
	@docker-compose -f docker/docker-compose.test.yaml down --volumes

##@ clean: Clean output files and build cache

.PHONY: clean

clean:
	@-echo "$(BOLD)$(RED)Removing build cache, test cache and files...$(RESET)"
	@-rm -rf $(PROJECT_BIN)
	go clean -testcache
	@cd $(SRC_STATIC_WEB) && \
  npm run clean && \
  cd -

##@ help: Help

.PHONY: help

help: Makefile
	@-echo " Usage:\n  make $(BLUE)<target>$(RESET)"
	@-echo
	@-sed -n 's/^##@//p' $< | column -t -s ':' | sed -e 's/[^ ]*/ &/2'
