.PHONY: all help list clean setup install build test unit integration coverage coverage-web run
.DEFAULT: all
.DEFAULT_GOAL: build

include Makefile.conf

GO := go
GO_OS ?= $(shell $(GO) env GOOS)
GO_ARCH ?= $(shell $(GO) env GOARCH)
GO_FLAGS ?= $(GO_FLAGS:)
GO_LINT := golint

# Get the current local branch name from git (if we can, this may be blank)
GIT_BRANCH := $(shell git symbolic-ref --short HEAD 2>/dev/null)
GIT_COMMIT := $(shell git rev-parse --short=7 HEAD 2>/dev/null)
# Get the git commit
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true 2>/dev/null)

# Build Flags
# The default version that's chosen when pushing the images. Can/should be overridden
BUILD_VERSION ?= $(shell git describe --abbrev=0 2>/dev/null)
BUILD_HASH = git-$(shell git rev-parse HEAD 2>/dev/null)
BUILD_TIME = $(shell date +%FT%T%z 2>/dev/null)

# If we don't set the build version it defaults to dev
ifeq ($(BUILD_VERSION),)
	BUILD_VERSION := dev
endif

BUILD_ENV =
ENV_FLAGS = CGO_ENABLED=1 $(BUILD_ENV)

GO_BUILD_FLAGS ?= -i -a -installsuffix cgo

# -tags netgo for enforcing the native Go DNS resolver
TAGS ?= -tags netgo

ifneq ($(GOOS), darwin)
	EXTLD_FLAGS = -extldflags "-lm -lstdc++ -static -v"
else
	EXTLD_FLAGS =
endif

# The flags we are passing to go build. -extldflags -static for making a static binary,
# -linkmode external for linking external C libraries into the binary, -X main.version for telling the
# Go binary which version it is.
GO_LINKER_FLAGS ?=-ldflags \
	'-s \
	-w \
	$(EXTLD_FLAGS) \
	-X "main.version=$(BUILD_VERSION)" \
	-X "${PACKAGE_BASE}/version.BuildHash=$(BUILD_HASH)" \
    -X "${PACKAGE_BASE}/core.BuildTime=$(BUILD_TIME)" \
    -X "${PACKAGE_BASE}/core.GitCommit=$(GIT_BRANCH)+$(GIT_COMMIT)$(GIT_DIRTY)"'

BUILD_FLAGS ?=$(GO_BUILD_FLAGS) $(TAGS) $(GO_LINKER_FLAGS)

PACKAGE_BASE = $(shell $(GO) list -e ./)
PKGS = $(shell $(GO) list . | grep -v /vendor/)

BINARY_BASE := $(BINARY_PATH)/$(shell basename `pwd`)
# Binary output name.
TARGET_BINARY := $(BINARY_BASE)-$(GO_OS)-$(GO_ARCH)

# To disable root, you can do "make SUDO="
SUDO := $(shell echo "sudo -E")
DOCKER := $(SUDO) docker

all: setup generate install test coverage-web verify format clean nuke build

# List all targets in this file
list:
	@$(MAKE) -rRpqn | awk -F':' '/^[a-z0-9][^$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)printf "$(HELP_COLOR)%-30s$(NO_COLOR)\n", A[i]}' | sort -u

help:
	@echo "Please use \`make <target>\`, Available options for <target> are:"
	@echo "$(HELP_COLOR)"
	@echo "  all                     to run all targets."
	@echo "  bench                   to run benchmark tests."
	@echo "  build                   to build out a binary."
	@echo "  clean                   to clean generated files."
	@echo "  clean-bin               to clean generated binaries only."
	@echo "  coverage                to run test with coverage."
	@echo "  coverage-web            to export coverage results to a web view \"coverage/profile.html\"."
	@echo "  deploy                  to deploy a docker container."
	@echo "  format                  to format the source code."
	@echo "  generate                to generate related files."
	@echo "  help                    to get help about the targets."
	@echo "  image                   to build a docker image."
	@echo "  install                 to install project dependencies."
	@echo "  integration             to run integration tests."
	@echo "  lint                    to run linter against source files."
	@echo "  list                    to list all targets."
	@echo "  nuke                    to enforce removing the corresponding installed archive or binary."
	@echo "  publish                 to publish the docker image."
	@echo "  run                     to run the generated binary, and build a new one if not existed."
	@echo "  setup                   to setup the tools used."
	@echo "  test                    to run all tests."
	@echo "  unit                    to run unit tests."
	@echo "  unit-with-race-cover    to run unit tests with race conditions coverage."
	@echo "  verify                  to lint and vet together."
	@echo "  vet                     to run detection on dead code."
	@echo "$(NO_COLOR)"


generate:
	@echo "$(OK_COLOR)==> Generating files via go generate...$(NO_COLOR)"
	@$(GO) generate $(GO_FLAGS) $(PKGS)

setup:
	@echo "$(OK_COLOR)==> Installing required components...$(NO_COLOR)"
	@$(GO) get -u $(GO_FLAGS) golang.org/x/tools/cmd/cover \
	github.com/golang/lint/golint

install:
	@echo "$(OK_COLOR)==> Installing packages into GOPATH...$(NO_COLOR)"
	@$(GO) install $(GO_FLAGS) $(TAGS) $(PKGS)

build: install \
    verify
	@echo "$(WARN_COLOR)==> Building binary...$(NO_COLOR)"
	for GO_OS in $(OSs); do \
		for GO_ARCH in $(ARCHS); do \
		    TARGET_BINARY=$(BINARY_BASE)-$$GO_OS-$$GO_ARCH; \
        	echo "$(INFO_COLOR) $(ENV_FLAGS) $(GO) build -o $$TARGET_BINARY $(BUILD_FLAGS) . $(NO_COLOR)\n"; \
        	env GOOS=$$GO_OS GOARCH=$$GO_ARCH \
            $(ENV_FLAGS) $(GO) build -o $$TARGET_BINARY $(BUILD_FLAGS) .; \
		done; \
	done

test: unit unit-with-race-cover bench integration coverage

# Unit tests
unit:
	@echo "$(WARN_COLOR)==> Unit tests...$(NO_COLOR)"
	$(GO) test -cover $(GO_FLAGS)  -timeout=8m ./...

unit-with-race-cover:
	@echo "$(WARN_COLOR)==> Unit tests with race cover...$(NO_COLOR)"
	@$(GO) test -race -cpu=1,2,4 -timeout 8m $(GO_FLAGS) ./...

bench:
	@echo "$(WARN_COLOR)==> Benchmarking tests...$(NO_COLOR)"
	@$(GO) test -run NONE -bench . -benchmem -tags 'bench' $(PKGS)

# Integration tests
integration: build verify
	@echo "$(WARN_COLOR)==> Integration tests...$(NO_COLOR)"
	$(GO) test -cover --tags=integration $(GO_FLAGS) ./...

# Coverage
coverage:
	@echo "$(WARN_COLOR)==> Coverage check...$(NO_COLOR)"
	@if [ ! -d $(COVERAGE_PATH) ] ; then mkdir -p $(COVERAGE_PATH) ; fi
	@$(GO) test -covermode=count -coverprofile $(COVERAGE_PATH)/profile.out $(GO_FLAGS) ./app

# Coverage using web view.
coverage-web:
	@echo "$(WARN_COLOR)==> Coverage web view export...$(NO_COLOR)"
	@if [ ! -d $(COVERAGE_PATH) ] ; then $(MAKE) $(COVERAGE_PATH) ; fi
	$(GO) tool cover -html=coverage/profile.out -o coverage/profile.html

#@if [ ! -d $(COVERAGE_PATH) ] ; then mkdir -p $(COVERAGE_PATH) ; fi
verify: vet lint

# Simplified dead code detector. Used for skipping certain checks on unreachable code
# (for instance, shift checks on arch-specific code).
# https://golang.org/cmd/vet/
vet:
	@echo "$(WARN_COLOR)==> Running vet...$(NO_COLOR)"
	@$(GO) vet $(GO_FLAGS) $(PKGS)

lint:
	@echo "$(WARN_COLOR)==> Running linter...$(NO_COLOR)"
	@$(GO_LINT) $(GO_FLAGS) $(PKGS)

format:
	@echo "$(WARN_COLOR)==> Formatting Code...$(NO_COLOR)"
	@$(GO) fmt $(GO_FLAGS) $(PKGS)

# Cleaning the project, by deleting binaries.
clean-bin:
	@echo "$(WARN_COLOR)==> Cleaning...$(NO_COLOR)"
	for GO_OS in $(OSs); do \
        for GO_ARCH in $(ARCHS); do \
            TARGET_BINARY=${BINARY_BASE}-$$GO_OS-$$GO_ARCH; \
            if [ -f $$TARGET_BINARY ] ; then rm -rf $$TARGET_BINARY ; fi \
        done; \
    done

clean: clean-bin
	$(GO) clean -i net
	rm -rf $(COVERAGE_PATH)

nuke:
	$(GO) clean -i ./...

run:
	if [ ! -f $(TARGET_BINARY) ] ; then $(MAKE) build; fi
	@$(TARGET_BINARY) $(args)

image:
	@echo "$(WARN_COLOR)==> Creating Docker Image...$(NO_COLOR)"
	@$(DOCKER) build . -t $(REGISTRY_REPO) $(args)

deploy:
	@echo "$(WARN_COLOR)==> Deploying Docker Container...$(NO_COLOR)"
	@$(SUDO) bash ./deploy.sh $(args)

publish:
	@echo "$(WARN_COLOR)==> Pushing Docker Image to $(REGISTRY_REPO)...$(NO_COLOR)"
	@$(DOCKER) push $(REGISTRY_REPO)
