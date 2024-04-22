PROJECT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOTESTSUM_VERSION := v1.11.0
GOFUMPT_VERSION := v0.6.0
RESPONSES_FILE ?= examples/example.json

bin/golangci-lint:
	@mkdir -p $(@D)
	GOBIN=$(PROJECT_DIR)/$(@D) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: bin/golangci-lint
	@bin/golangci-lint run

bin/gofumpt:
	@mkdir -p $(@D)
	GOBIN=$(PROJECT_DIR)/$(@D) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

format: bin/gofumpt
	@bin/gofumpt -w .

bin/gotestsum:
	@mkdir -p $(@D)
	GOBIN=$(PROJECT_DIR)/$(@D) go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)

test: bin/gotestsum lint
	@cd $(PROJECT_DIR) && $(PROJECT_DIR)/bin/gotestsum --format testdox -- -coverprofile=cover.out ./internal/...
	@cd $(PROJECT_DIR) && go tool cover -func=cover.out > coverage-text.txt

build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/pretender cmd/pretender/main.go

run:
	go run cmd/pretender/main.go --responses $(RESPONSES_FILE)

docker-build:
	docker build $(PROJECT_DIR) -t pretender:latest

docker-run: docker-build
	docker run --rm -v $(PROJECT_DIR)/$(RESPONSES_FILE):/$(RESPONSES_FILE) -p 8080:8080 pretender:latest --responses /$(RESPONSES_FILE)

version-check:
	@go run tools/versioncheck/main.go

.PHONY: lint format test build run docker-build docker-run version-check
