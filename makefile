APP_NAME		= go-rest-api-template

SHELL			= bash

GO := $(shell command -v go 2> /dev/null)
ifndef GO
$(error go is required, please install)
endif

# Dependencies
GOPATH			:= $(shell go env GOPATH)
GOBIN			?= $(GOPATH)/bin
GOFUMPT			:= $(GOBIN)/gofumpt
GOLANGCILINT   	:= $(GOBIN)/golangci-lint

# Paths
FILES = $(shell find . -name '.?*' -prune -o -name vendor -prune -o -name '*.go' -print)

PKGS  = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f \
            '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
            $(PKGS))

# Directories
BUILD_DIR     := build

.PHONY: help
help: ## Print this menu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-gofumpt: ## Install Go formatter.
	$(GO) get -u mvdan.cc/gofumpt

fmt: install-gofumpt  ## Format code
	$(GO) fmt $(PKGS)
	$(GOFUMPT) -s -w $(FILES)

install-golangcilint: ## Install golangcilint
	# To bump, simply change the version at the end to the desired version. The git sha here points to the newest commit
	# of the install script verified by our team located here: https://github.com/golangci/golangci-lint/blob/master/install.sh
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/17d24ebd671875cdf52804e1ca72ca8f0718a844/install.sh | sh -s -- -b ${GOBIN} v1.40.1

lint: install-golangcilint ## Lint
	$(GOLANGCILINT) run

check: fmt lint test ## Run this continously to verify that all CI checks are green

test: ## Test app
	go test $(TESTPKGS)

build: ## Build app
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)

run: ## Run app
	go run *.go

build-docker: ## Build docker image
	docker build -t $(APP_NAME) .

run-docker: build-docker ## Run app with docker
	docker run \
		 --name $(APP_NAME) \
		 --rm \
		 -e PORT="8080" \
		 -p 8080:8080 \
		 $(APP_NAME)

# Coverage
GOCOVMERGE      := $(GOBIN)/gocovmerge
GOCOVXML        := $(GOBIN)/gocov-xml
GOCOV           := $(GOBIN)/gocov

COVERAGE_MODE    = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML     = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML    = $(COVERAGE_DIR)/index.html

$(GOCOVMERGE):
	$(GO) install github.com/wadey/gocovmerge@latest

$(GOCOVXML):
	$(GO) install github.com/AlekSi/gocov-xml@latest

$(GOCOV):
	$(GO) install github.com/axw/gocov/gocov@v1.0.0

test-coverage-tools: | $(GOCOVMERGE) $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(BUILD_DIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint test-coverage-tools
	@mkdir -p $(COVERAGE_DIR)/coverage
	@for pkg in $(TESTPKGS); do \
        go test \
            -coverpkg=$$(go list -f '{{ join .Deps "\n" }}' $$pkg | \
                    grep '^$(MODULE)/' | \
                    tr '\n' ',')$$pkg \
            -covermode=$(COVERAGE_MODE) \
            -coverprofile="$(COVERAGE_DIR)/coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
     done
	@$(GOCOVMERGE) $(COVERAGE_DIR)/coverage/*.cover > $(COVERAGE_PROFILE)
	@$(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	@$(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)
