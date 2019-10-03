PREFIX?=$(shell pwd)
NAME := $(shell cat package.json | jq -r .name)
PKG := github.com/samygp/$(NAME)
BUILDDIR := ${PREFIX}/dist
ENTRYPOINT := cmd/main.go
VERSION := $(shell cat package.json | jq -r .version)
CTIMEVAR=-X $(PKG)/version.Name=$(NAME) -X $(PKG)/version.Version=$(VERSION)
GO_LDFLAGS=-ldflags "$(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"
SEMBUMP_IMG=chatu/sembump:0.1.0
SEMBUMP=docker run --rm -v `pwd`:/app $(SEMBUMP_IMG)

.PHONY: default
default: help

.PHONY: name
name: ## Output name of project
	@echo $(NAME)

.PHONY: version
version: ## Output current version
	@echo $(VERSION)

.PHONY: install-deps
install-deps: ## Install dependencies
	@echo "+ $@"
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/kisielk/errcheck
	@go get -u honnef.co/go/tools/cmd/staticcheck

.PHONY: local-build
local-build: ## Builds a dynamic executable or package
	@echo "+ $@"
	@go build \
	  ${GO_LDFLAGS} -o $(NAME) \
	  $(ENTRYPOINT)

.PHONY: pi-build
pi-build: ## Builds a dynamic executable or package
	@echo "+ $@"
	GOOS=linux GOARCH=arm GOARM=5 go build \
	  -o pi-$(NAME) \
	  $(ENTRYPOINT)

.PROXY: run
run: ## Execute built binary as web
	@echo "+ $@"
	@$(shell grep -v ^# config.env | xargs) ./$(NAME)

.PHONY: runmon
runmon: restart ## Run executable and restart it if any changes are detected
	@fswatch -o $(shell find . -type f -name '*.go' -not -path "./vendor/*") config.env | xargs -n1 -I{} make restart || make kill

.PHONY: restart
restart: kill local-build
	@echo "+ $@"
	@$(shell grep -v ^# config.env | xargs) ./$(NAME) &

.PHONY: kill
kill:
	@echo "+ $@"
	@pgrep $(NAME) | xargs kill || true

.PROXY: static
static: ## Generate static binary
	@echo "+ $@"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
	  -o $(NAME) \
	  -a -tags "static_build netgo" \
	  -installsuffix netgo ${GO_LDFLAGS_STATIC} \
	  $(ENTRYPOINT)

.PHONY: docker-build
docker-build: ## Build docker image
	@echo "+ $@"
	@docker $$DOCKER_CONFIG build -t $(NAME) .

.PHONY: docker-run
docker-run: #docker-build ## Execute built docker image
	@echo "+ $@"
	@docker run -it --privileged --rm \
	 -h $$HOSTNAME \
	 -t --env-file $(shell pwd)/config.env \
	 $(NAME)

.PHONY: verify
verify: fmt lint vet errcheck staticcheck ## Verify code for common issues

.PHONY: fmt
fmt: ## Verifies all files have been `gofmt`ed
	@echo "+ $@"
	@test -z "$$(gofmt -s -l $(shell find . -type f -name '*.go' -not -path "./vendor/*") 2>&1 | tee /dev/stderr)"

.PHONY: lint
lint: ## Verifies `golint` passes
	@echo "+ $@"
	@test -z "$$(golint $(shell go list ./... | grep -v vendor) 2>&1 | tee /dev/stderr)"

.PHONY: clean
clean: ## Cleanup any build binaries or packages
	@echo "+ $@"
	$(RM) $(NAME)
	$(RM) -r $(BUILDDIR)

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
