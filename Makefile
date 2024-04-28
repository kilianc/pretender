PROJECT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOTESTSUM_VERSION := v1.11.0
GOFUMPT_VERSION := v0.6.0
GOLANGCI_LINT_VERSION := v1.57.2
RESPONSES_FILE ?= examples/example.json
BINARY_NAME := pretender
OS_LIST := darwin linux
ARCH_LIST := arm64 amd64
BUILD_TARGETS := $(foreach os,$(OS_LIST),$(foreach arch,$(ARCH_LIST),bin/$(BINARY_NAME)-$(os)-$(arch)))
RELEASE_TARGETS := $(foreach os,$(OS_LIST),$(foreach arch,$(ARCH_LIST),bin/$(BINARY_NAME)-$(os)-$(arch).tar.gz))

# - install binary dependencies

bin/golangci-lint:
	@mkdir -p $(@D)
	GOBIN=$(PROJECT_DIR)/$(@D) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

bin/gofumpt:
	@mkdir -p $(@D)
	GOBIN=$(PROJECT_DIR)/$(@D) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

bin/gotestsum:
	@mkdir -p $(@D)
	GOBIN=$(PROJECT_DIR)/$(@D) go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)

bin/git-chglog:
	@mkdir -p $(@D)
	GOBIN=$(PROJECT_DIR)/$(@D) go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

# - run, test, lint, format etc.

.PHONY: run
run:
	go run cmd/$(BINARY_NAME)/main.go --responses $(RESPONSES_FILE)

.PHONY: lint
lint: bin/golangci-lint
	@bin/golangci-lint run

.PHONY: format
format: bin/gofumpt
	@bin/gofumpt -w .

.PHONY: test
test: bin/gotestsum lint
	@echo ""
	@cd $(PROJECT_DIR) && $(PROJECT_DIR)/bin/gotestsum --format testdox -- -coverprofile=cover.out -coverpkg=./... $(shell go list ./... | grep -v /tools/ | grep -v /cmd/)

cover.out:
	@if [ ! -f cover.out ]; then $(MAKE) test; fi

cover.txt: cover.out
	@cd $(PROJECT_DIR) && go tool cover -func=cover.out -o cover.txt

.PHONY: open-cover
open-cover: cover.out
	@cd $(PROJECT_DIR) && go tool cover -html=cover.out

# - build and release

.PHONY: build
build: $(BUILD_TARGETS)

.PHONY: $(BUILD_TARGETS)
$(BUILD_TARGETS):
	@$(eval os = $(word 2, $(subst -, ,$@)))
	@$(eval arch = $(word 3, $(subst -, ,$@)))
	GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=0 go build -ldflags "-s -w" -o $@ cmd/$(BINARY_NAME)/main.go

.PHONY: release
release: VERSION = $(shell go run cmd/$(BINARY_NAME)/main.go --responses $(RESPONSES_FILE) --version)
release: $(RELEASE_TARGETS)
	GOPROXY=proxy.golang.org go list -m github.com/kilianc/$(BINARY_NAME)@$(VERSION)

.PHONY: $(RELEASE_TARGETS)
$(RELEASE_TARGETS): clean build
	@cp $(shell echo $@ | sed s/.tar.gz//) bin/$(BINARY_NAME)
	cd bin && tar -czf $(shell basename $@) $(BINARY_NAME)
	@rm bin/$(BINARY_NAME)

.PHONY: docker-build
docker-build:
	docker build $(PROJECT_DIR) -t $(BINARY_NAME):latest

.PHONY: docker-run
docker-run: docker-build
	docker run --rm -v $(PROJECT_DIR)/$(RESPONSES_FILE):/$(RESPONSES_FILE) -p 8080:8080 $(BINARY_NAME):latest --responses /$(RESPONSES_FILE)

# - tools

.PHONY: changelog
changelog: bin/git-chglog
	@bin/git-chglog -o CHANGELOG.md --next-tag $(tag)

.PHONY: version-check
version-check:
	@go run tools/versioncheck/main.go $(tag)

.PHONY: cover-check
cover-check: cover.txt
	@echo ""
	@cat cover.txt
	@go run tools/covercheck/main.go

.PHONY: commit-check
commit-check:
	@echo ""
	@go run tools/commitcheck/main.go '$(message)'

# - clean

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf cover.*
