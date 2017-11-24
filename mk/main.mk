GO := go
GO_OS ?= $(shell $(GO) env GOOS)
GO_ARCH ?= $(shell $(GO) env GOARCH)
GO_FLAGS ?= $(GO_FLAGS:)
GO_LINT ?= golint
GO_FMT ?= gofmt

PKG_BASE ?= $(shell $(GO) list -e ./)
PKGS ?= $(shell $(GO) list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")

# To disable root, you can do "make SUDO="
SUDO := $(shell echo "sudo -E")
DOCKER := $(shell docker info > /dev/null 2>&1 || $(SUDO)) docker

include mk/build.mk
include mk/coverage.mk
include mk/docker.mk
include mk/test.mk
include mk/validate.mk

TARGET_BINARY := $(BINARY_PATH)/$(BINARY_PREFIX)-$(GO_OS)-$(GO_ARCH)
