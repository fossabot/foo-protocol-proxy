.PHONY: all
.DEFAULT: all
.DEFAULT_GOAL: build

include Makefile.conf

GO := go
GO_OS ?= $(shell $(GO) env GOOS)
GO_ARCH ?= $(shell $(GO) env GOARCH)
GO_FLAGS ?= $(GO_FLAGS:)
GO_LINT := golint

# Get the current local branch name from git (if we can, this may be blank)
GIT_BRANCH := $(shell git symbolic-ref --short HEAD || git symbolic-ref --short HEAD 2> /dev/null)
GIT_COMMIT := $(shell git rev-parse HEAD 2> /dev/null)
# Get the git commit
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true 2> /dev/null)

# Build Flags
# The default version that's chosen when pushing the images. Can/should be overridden
BUILD_VERSION ?= $(shell git describe --abbrev=0 | cut -d "v" -f 2 2> /dev/null)
BUILD_HASH ?= git-$(shell git rev-parse --short=18 HEAD 2> /dev/null)
BUILD_TIME ?= $(shell date +%FT%T%z 2> /dev/null)

# If we don't set the build version it defaults to dev
ifeq ($(BUILD_VERSION),)
	BUILD_VERSION := $(shell cat $(CURDIR)/.version 2> /dev/null || echo dev)
endif

BUILD_ENV =

GO_BUILD_FLAGS ?= -i -a -installsuffix cgo

EXTLD_FLAGS =

ifneq ($(GOOS), darwin)
	EXTLD_FLAGS = -extldflags "-lm -lstdc++ -static -v"
endif

# The flags we are passing to go build. -extldflags -static for making a static binary,
# -linkmode external for linking external C libraries into the binary, -X main.version for telling the
# Go binary which version it is.
GO_LINKER_FLAGS ?=-s \
	-w \
	$(EXTLD_FLAGS) \
	-X "${PACKAGE_BASE}/version.BuildHash=$(BUILD_HASH)" \
	-X "${PACKAGE_BASE}/core.BuildTime=$(BUILD_TIME)" \
	-X "${PACKAGE_BASE}/core.GitBranch=$(GIT_BRANCH)" \
	-X "${PACKAGE_BASE}/core.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)"

ifdef BUILD_VERSION
	GO_LINKER_FLAGS += -X main.version=$(BUILD_VERSION)
	DOCKER_TAG = $(BUILD_VERSION)
endif

GO_ENV_FLAGS ?= CGO_ENABLED=1

# netgo for enforcing the native Go DNS resolver
GO_TAGS ?= netgo

BUILD_FLAGS ?= $(GO_BUILD_FLAGS) -ldflags '$(GO_LINKER_FLAGS)' -tags $(GO_TAGS) $(GO_FLAGS)

ifneq ($(CROSS_BUILD), true)
    PLATFORMS = $(GO_OS)
    ARCHCS = $(GO_ARCH)
    GO_ENV_FLAGS = CGO_ENABLED=0
endif

GO_ENV_FLAGS += $(BUILD_ENV)

PACKAGE_BASE = $(shell $(GO) list -e ./)
PKGS = $(shell $(GO) list ./... | grep -v /vendor/)

BINARY_BASE := $(BINARY_PATH)/$(shell basename `pwd`)

# Binary output name.
TARGET_BINARY := $(BINARY_BASE)-$(GO_OS)-$(GO_ARCH)

BUILD_TARGET_PREFIX := build-bin-for-
buildTargets = $(foreach a, $(3), $(foreach b, $(4), $(call $(1), $(2)$(a), -$(b))))
BUILD_TARGETS = $(call buildTargets, addprefix, $(BUILD_TARGET_PREFIX), $(PLATFORMS), $(ARCHS))

# To disable root, you can do "make SUDO="
SUDO := $(shell echo "sudo -E")
DOCKER := $(shell docker info > /dev/null 2>&1 || $(SUDO)) docker

all: setup generate install test coverage-web verify format clean nuke build deploy list help

# List all targets in this file
list:
	@$(MAKE) -rRpqn | awk -F':' '/^[a-z0-9][^$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)printf "$(DISCLAIMER_COLOR)%-30s$(NO_COLOR)\n", A[i]}' | sort -u

help:
	@echo "$(OK_COLOR)"
	@echo "$$FOO_PROTOCOL_PROXY"
	@echo "$(NO_COLOR)"
	@echo "Please use \`make <target>\`, Available options for <target> are:"
	@echo "$(HELP_COLOR)"
	@echo "  all                     to run all targets."
	@echo "  bench                   to run benchmark tests."
	@echo "  build-bin               to build out a binary."
	@echo "  build                   to install dependencies and build out a binary."
	@echo "  clean                   to clean generated files."
	@echo "  clean-bin               to clean generated binaries only."
	@echo "  cover                   to run test with coverage and report that out to profile."
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
	@echo "  setup                   to setup the external used tools."
	@echo "  test                    to run all tests."
	@echo "  unit                    to run unit tests."
	@echo "  unit-with-race-cover    to run unit tests with race conditions coverage."
	@echo "  verify                  to lint and vet together."
	@echo "  vet                     to run detection on dead code."
	@echo "$(NO_COLOR)"

setup:
	@echo "$(OK_COLOR)$(MSG_PREFIX) Setting-up required components...$(NO_COLOR)"
	$(GO) get -u $(GO_FLAGS) golang.org/x/tools/cmd/cover \
		github.com/golang/lint/golint

generate:
	@echo "$(OK_COLOR)$(MSG_PREFIX) Generating files via go generate...$(NO_COLOR)"
	@$(GO) generate $(GO_FLAGS) $(PKGS)

install:
	@echo "$(OK_COLOR)$(MSG_PREFIX) Installing packages into GOPATH...$(NO_COLOR)"
	@$(GO) install $(GO_FLAGS) -tags $(GO_TAGS) $(PKGS)

build-bin-for-%:
	@$(eval GO_OS=$(firstword $(subst -, , $*)))
	@$(eval GO_ARCH=$(or $(word 2,$(subst -, , $*)),$(value 2)))
	@$(eval TARGET_BINARY=$(BINARY_BASE)-$(GO_OS)-$(GO_ARCH))
	@echo "$(INFO_COLOR) GOOS=$(GO_OS) GOARCH=$(GO_ARCH) $(GO_ENV_FLAGS) $(GO) build -o $(TARGET_BINARY) $(BUILD_FLAGS) . $(NO_COLOR)\n"
	@env GOOS=$(GO_OS) GOARCH=$(GO_ARCH) $(GO_ENV_FLAGS) $(GO) build -o $(TARGET_BINARY) $(BUILD_FLAGS) .

build-bin:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Building binary...$(NO_COLOR)"
	@$(MAKE) $(BUILD_TARGETS)

ifneq ($(BUILD_IN_CONTAINER), true)

build: install build-bin

else

build: deploy

endif

test: unit unit-with-race-cover bench integration cover

# Unit tests
unit:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Unit tests...$(NO_COLOR)"
	@$(GO) test -short -cover -timeout=$(TEST_TIMEOUT) -tags $(GO_TAGS) $(GO_FLAGS) ./...

unit-with-race-cover:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Unit tests with race cover...$(NO_COLOR)"
	@$(GO) test -race -cpu=1,2,4 -timeout $(TEST_TIMEOUT) -tags $(GO_TAGS) $(GO_FLAGS) ./...

bench:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Benchmarking tests...$(NO_COLOR)"
	@$(GO) test -run NONE -bench . -benchmem -tags bench $(GO_FLAGS) $(PKGS)

# Integration tests
integration: build verify
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Integration tests...$(NO_COLOR)"
	@$(GO) test -cover -tags integration $(GO_FLAGS) ./...

# Coverage
cover:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Coverage check...$(NO_COLOR)"
	@if [ ! -d $(COVERAGE_PATH) ] ; then mkdir -p $(COVERAGE_PATH) ; fi
	@$(GO) test -covermode=$(COVERAGE_MODE) -coverprofile $(COVERAGE_PATH)/app.part $(GO_FLAGS) ./app
	@$(GO) test -covermode=$(COVERAGE_MODE) -coverprofile $(COVERAGE_PATH)/handlers.part $(GO_FLAGS) ./handlers
	@$(GO) test -covermode=$(COVERAGE_MODE) -coverprofile $(COVERAGE_PATH)/persistence.part $(GO_FLAGS) ./persistence
	@echo "mode: count" > $(COVERAGE_PROFILE)
	@grep -h -v "mode: count" $(COVERAGE_PATH)/*.part >> $(COVERAGE_PROFILE)

# Coverage using web view.
coverage-web:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Coverage web view export...$(NO_COLOR)"
	@if [ ! -d $(COVERAGE_PATH) ] ; then $(MAKE) $(COVERAGE_PATH) ; fi
	@$(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML) $(GO_FLAGS)

verify: vet lint

# Simplified dead code detector. Used for skipping certain checks on unreachable code
# (for instance, shift checks on arch-specific code).
# https://golang.org/cmd/vet/
vet:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Running vet...$(NO_COLOR)"
	@$(GO) vet $(GO_FLAGS) $(PKGS)

lint:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Running linter...$(NO_COLOR)"
	@$(GO_LINT) $(GO_FLAGS) $(PKGS)

format:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Formatting Code...$(NO_COLOR)"
	@$(GO) fmt $(GO_FLAGS) $(PKGS)

# Deleting binaries.
clean-bin:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Cleaning binaries...$(NO_COLOR)"
	for GO_OS in $(PLATFORMS); do \
		for GO_ARCH in $(ARCHS); do \
			TARGET_BINARY=${BINARY_BASE}-$$GO_OS-$$GO_ARCH; \
			if [ -f $$TARGET_BINARY ] ; then rm -rf $$TARGET_BINARY ; fi \
		done; \
	done

# Cleaning up.
clean: clean-bin
	@$(GO) clean -i $(GO_FLAGS) net
	@rm -rf $(COVERAGE_PATH)

nuke:
	@$(GO) clean -i $(GO_FLAGS) ./...

run:
	@if [ ! -f $(TARGET_BINARY) ] ; then $(MAKE) build; fi
	@$(TARGET_BINARY) $(args)

image:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Creating Docker Image...$(NO_COLOR)"
	@$(DOCKER) build . -t $(REGISTRY_REPO):$(DOCKER_TAG) -f $(DOCKER_FILE) $(args)

deploy:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Deploying Docker Container...$(NO_COLOR)"
	@$(SUDO) bash ./deploy.sh $(args)

publish:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Pushing Docker Image to $(REGISTRY_REPO)...$(NO_COLOR)"
	@$(DOCKER) push $(REGISTRY_REPO):$(DOCKER_TAG)
