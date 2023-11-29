####
## Runtime variables
####

ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN_DIR := $(ROOT_DIR)/bin
GOMODULES := $(shell find $(ROOT_DIR) -name 'go.mod')
BUILD_TIMESTAMP := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_HASH := $(shell git rev-parse HEAD)

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

.PHONY: test
test: ## Run Go unit tests
	@go test ./...

.PHONY: gomod-tidy
gomod-tidy:
	go mod tidy

.PHONY: license-cache
license-cache: bin/licensei ## Generate license cache
	$(LICENSEI_BIN) cache

.PHONY: license-check
license-check: bin/licensei license-cache ## Check licenses for software components
	$(LICENSEI_BIN) check

.PHONY: license-header
license-header: bin/licensei ## Check license headers in source code files
	$(LICENSEI_BIN) header

.PHONY: lint-actions
lint-actions: bin/actionlint ## Lint Github Actions
	@$(ACTIONLINT_BIN) -color

LINTGOMODULES = $(addprefix lint-, $(GOMODULES))

.PHONY: $(LINTGOMODULES)
$(LINTGOMODULES):
	cd $(dir $(@:lint-%=%)) && "$(GOLANGCI_BIN)" run -c "$(GOLANGCI_CONFIG)"

.PHONY: lint-go
lint-go: bin/golangci-lint $(LINTGOMODULES) ## Lint Go source code

####
##  Golangci-lint CLI
####

GOLANGCI_BIN := $(BIN_DIR)/golangci-lint
GOLANGCI_CONFIG := $(ROOT_DIR)/.golangci.yml
GOLANGCI_VERSION := 1.54.2

bin/golangci-lint: bin/golangci-lint-$(GOLANGCI_VERSION)
	@ln -sf golangci-lint-$(GOLANGCI_VERSION) bin/golangci-lint

bin/golangci-lint-$(GOLANGCI_VERSION): | $(BIN_DIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b "$(BIN_DIR)" "v$(GOLANGCI_VERSION)"
	@mv bin/golangci-lint $@

####
## ActionLint CLI
####

ACTIONLINT_BIN := $(BIN_DIR)/actionlint
ACTIONLINT_VERSION := 1.6.26

bin/actionlint: bin/actionlint-$(ACTIONLINT_VERSION)
	@ln -sf actionlint-$(ACTIONLINT_VERSION) bin/actionlint

bin/actionlint-$(ACTIONLINT_VERSION): | $(BIN_DIR)
	curl -sSfL https://raw.githubusercontent.com/rhysd/actionlint/main/scripts/download-actionlint.bash \
	| bash -s -- "$(ACTIONLINT_VERSION)" "$(BIN_DIR)"
	@mv bin/actionlint $@

####
## Licensei CLI
####

LICENSEI_BIN := $(BIN_DIR)/licensei
LICENSEI_VERSION = 0.9.0

bin/licensei: bin/licensei-$(LICENSEI_VERSION)
	@ln -sf licensei-$(LICENSEI_VERSION) bin/licensei

bin/licensei-$(LICENSEI_VERSION): | $(BIN_DIR)
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s v$(LICENSEI_VERSION)
	@mv bin/licensei $@
