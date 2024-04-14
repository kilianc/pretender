PROJECT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOTESTSUM_VERSION := v1.11.0
RESPONSES ?= README.md

bin/gotestsum:
	@mkdir -p $(@D)
	GOBIN=$(PROJECT_DIR)/$(@D) go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)

test: bin/gotestsum
	cd $(PROJECT_DIR) && $(PROJECT_DIR)/bin/gotestsum --format testdox

build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/pretender cmd/pretender/main.go

run:
	go run cmd/pretender/main.go --responses $(RESPONSES)

docker-build:
	docker build $(PROJECT_DIR) -t pretender:latest

docker-run: docker-build
	docker run --rm -v $(PROJECT_DIR)/README.md:/README.md -p 8080:8080 pretender:latest --responses /README.md

version-check:
	go run scripts/versioncheck.go

.PHONY: test build run docker-build docker-run version-check
