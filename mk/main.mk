GO := go
GO_OS ?= $(shell $(GO) env GOOS 2> /dev/null)
GO_ARCH ?= $(shell $(GO) env GOARCH 2> /dev/null)
GO_FLAGS ?= $(GO_FLAGS:)
GO_LINT ?= golint
GO_FMT ?= gofmt

PKG_BASE ?= $(shell $(GO) list -e ./ 2> /dev/null)
PKGS ?= $(shell $(GO) list ./... 2> /dev/null | grep -v /vendor/)
GO_FILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*" 2> /dev/null)

include mk/utils.mk
include mk/build.mk
include mk/coverage.mk
include mk/docker.mk
include mk/test.mk
include mk/validate.mk

TARGET_BINARY := $(BINARY_PATH)/$(BINARY_PREFIX)-$(GO_OS)-$(GO_ARCH)
