# Image URL to use all building/pushing image targets
IMG ?= nx2internal.azurecr.io/timref/matrix-go-push:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: docker-build

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: test
test: fmt vet staticcheck
	go test ./... -coverprofile cover.out

.PHONY: docker-build
docker-build: test
	docker build -t ${IMG} .

.PHONY: docker-build-gitlabci
docker-build-gitlabci:
	docker build -t ${IMG} .

.PHONY: run
run: fmt vet
	go run ./cmd/main.go