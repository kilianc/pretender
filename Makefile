PROJECT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOTESTSUM_VERSION := v1.11.0
GOFUMPT_VERSION := v0.6.0
GOLANGCI_LINT_VERSION := v1.57.2
RESPONSES_FILE ?= examples/example.json
BINARY_NAME := pretender
OS_LIST := darwin linux
ARCH_LIST := arm64 amd64
BUILD_TARGETS := $(foreach os,$(OS_LIST),$(foreach arch,$(ARCH_LIST),bin/$(BINARY_NAME)_$(os)_$(arch)))
RELEASE_TARGETS := $(foreach os,$(OS_LIST),$(foreach arch,$(ARCH_LIST),bin/$(BINARY_NAME)_$(os)_$(arch).tar.gz))

bin/golangci-lint:
	@mkdir -p $(@D)
	GOBIN=$(PROJECT_DIR)/$(@D) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

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

cover: test
	@cd $(PROJECT_DIR) && go tool cover -html=cover.out

build: $(BUILD_TARGETS)

$(BUILD_TARGETS):
	@$(eval os = $(word 2, $(subst _, ,$@)))
	@$(eval arch = $(word 3, $(subst _, ,$@)))
	GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=0 go build -ldflags "-s -w" -o $@ cmd/$(BINARY_NAME)/main.go

release: VERSION = $(shell go run cmd/$(BINARY_NAME)/main.go --responses $(RESPONSES_FILE) --version)
release: $(RELEASE_TARGETS)
	GOPROXY=proxy.golang.org go list -m github.com/kilianc/$(BINARY_NAME)@$(VERSION)

$(RELEASE_TARGETS): clean build
	@cp $(shell echo $@ | sed s/.tar.gz//) bin/$(BINARY_NAME)
	cd bin && tar -czf $(shell basename $@) $(BINARY_NAME)
	@rm bin/$(BINARY_NAME)

bin/git-chglog:
	@mkdir -p $(@D)
	GOBIN=$(PROJECT_DIR)/$(@D) go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

changelog: bin/git-chglog
	@bin/git-chglog -o CHANGELOG.md --next-tag $(tag)

run:
	go run cmd/$(BINARY_NAME)/main.go --responses $(RESPONSES_FILE)

docker-build:
	docker build $(PROJECT_DIR) -t $(BINARY_NAME):latest

docker-run: docker-build
	docker run --rm -v $(PROJECT_DIR)/$(RESPONSES_FILE):/$(RESPONSES_FILE) -p 8080:8080 $(BINARY_NAME):latest --responses /$(RESPONSES_FILE)

version-check:
	@go run tools/versioncheck/main.go $(tag)

clean:
	rm -rf bin/*

.PHONY: lint format test build release run docker-build docker-run version-check clean
