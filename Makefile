# ==============================================================================
# onexstack protoc plugin collection.
#
# Each plugin under cmd/ is built and installed on its own. A plugin whose
# directory is named protoc-gen-X is invoked by protoc as --X_out, so the
# directory name is not cosmetic: renaming one renames the flag its consumers
# must pass.

SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
PROJ_ROOT_DIR := $(strip $(abspath $(shell cd $(SELF_DIR)/ && pwd -P)))
OUTPUT_DIR := $(PROJ_ROOT_DIR)/_output

# Every plugin builds today. The two that did not are worth recording, because
# both failed the same way -- an import path left over from wherever the source
# was copied from, pointing at something that does not exist:
#
#   protoc-gen-go-errordoc  imported
#     github.com/onexstack/onex/tools/protoc-gen-go-errors-code/errors while
#     shipping a local errors/ package beside it. It is now
#     protoc-gen-go-errors-code, so the directory matches the --go-errors-code_out
#     flag its own template documents.
#   protoc-gen-go-defaults  is whole-namespace consistent now; its earlier state
#     is described in its README.
#
# A leftover import is not caught by building the plugin you are working on: it
# only surfaces when the whole module is loaded, which is what `go mod tidy`
# does. That is why this list, rather than the imports, is what went stale.
#
# protoc-gen-go-defaults pulls in ./defaults as well: that package is the
# runtime half of the plugin -- the generated code blank-imports it to register
# the (defaults.value) extension -- so it is built and vetted with the plugin
# even though it lives outside cmd/.
PLUGINS ?= ./defaults/... ./cmd/protoc-gen-go-deepcopy/... ./cmd/protoc-gen-go-defaults/... ./cmd/protoc-gen-go-errors-code/...

.PHONY: all
all: build vet test

.PHONY: build
build: ## Compile every plugin that currently builds.
	@echo "===========> Building plugins"
	@mkdir -p $(OUTPUT_DIR)
	@go build -o $(OUTPUT_DIR)/ $(PLUGINS)

.PHONY: vet
vet: ## Run go vet over the plugins.
	@echo "===========> Running go vet"
	@go vet $(PLUGINS)

.PHONY: test
test: ## Run the plugins' unit tests.
	@echo "===========> Running unit tests"
	@go test -race -count=1 $(PLUGINS)

.PHONY: fmt
fmt: ## Format Go sources with gofmt (-s -w).
	@echo "===========> Formatting sources"
	@gofmt -s -w ./cmd

.PHONY: install
install: ## Install the plugins into $(go env GOPATH)/bin for protoc to find.
	@echo "===========> Installing protoc plugins"
	@go install $(PLUGINS)

.PHONY: clean
clean: ## Remove build artifacts (_output/).
	@-rm -vrf $(OUTPUT_DIR)

help: Makefile ## Show available targets.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<TARGET>\033[0m\n\n\033[35mTargets:\033[0m\n"} /^[0-9A-Za-z._-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' Makefile
