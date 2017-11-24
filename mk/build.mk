# Get the current local branch name from git (if we can, this may be blank)
GIT_BRANCH := $(shell git symbolic-ref --short HEAD || git symbolic-ref --short HEAD 2> /dev/null)
GIT_COMMIT := $(shell git rev-parse HEAD 2> /dev/null)
# Get the git commit
GIT_DIRTY := $(shell test -n "`git status --porcelain --untracked-files=no`" && echo "+CHANGES" || true 2> /dev/null)

# Build Flags
# The default version that's chosen when pushing the images. Can/should be overridden
BUILD_VERSION ?= $(shell git describe --abbrev=0 | cut -d "v" -f 2 2> /dev/null)
BUILD_HASH ?= git-$(shell git rev-parse --short=18 HEAD 2> /dev/null)
BUILD_TIME ?= $(shell date +%FT%T%z 2> /dev/null)

# If we don't set the build version it defaults to dev
ifeq ($(BUILD_VERSION),)
	BUILD_VERSION := $(shell cat $(CURDIR)/.version 2> /dev/null || echo dev)
endif

BUILD_ENV ?= $(BUILD_ENV:)

GO_BUILD_FLAGS ?= -i -a -installsuffix cgo

EXTLD_FLAGS ?=

# Honor static
ifeq ($(STATIC),true)
	# Append to the version
	EXTLD_FLAGS += -static
endif

ifneq ($(GOOS), darwin)
	EXTLD_FLAGS := "-lm -lstdc++ -v"
endif

# The flags we are passing to go build. -extldflags -static for making a static binary,
# -linkmode external for linking external C libraries into the binary, -X main.version for telling the
# Go binary which version it is.
GO_LINKER_FLAGS ?=-s \
	-w \
	-extldflags $(EXTLD_FLAGS) \
	-X "${PKG_BASE}/version.BuildHash=$(BUILD_HASH)" \
	-X "${PKG_BASE}/core.BuildTime=$(BUILD_TIME)" \
	-X "${PKG_BASE}/core.GitBranch=$(GIT_BRANCH)" \
	-X "${PKG_BASE}/core.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)"

ifdef BUILD_VERSION
	GO_LINKER_FLAGS += -X main.version=$(BUILD_VERSION)
	DOCKER_TAG = $(BUILD_VERSION)
endif

# Honor debug
ifeq ($(DEBUG),true)
	# Disable function inlining and variable registration
	GO_GC_FLAGS := -gcflags "-N -l"
endif

# netgo for enforcing the native Go DNS resolver
GO_TAGS ?= netgo

BUILD_FLAGS ?= $(GO_BUILD_FLAGS) -ldflags '$(GO_LINKER_FLAGS)' $(GO_GC_FLAGS) -tags $(GO_TAGS) $(GO_FLAGS)

# Binary output prefix.
BINARY_PREFIX := $(shell basename `pwd`)

GO_ENV_FLAGS ?= $(GO_ENV_FLAGS:)
GO_ENV_FLAGS += $(BUILD_ENV)

extension = $(patsubst windows, .exe, $(filter windows, $(1)))

# Valid target combinations
VALID_OS_ARCH := "[darwin/amd64][linux/amd64][linux/arm][linux/arm64][windows/amd64][windows/386]"

os.darwin := darwin
os.linux := linux
os.windows := windows

arch.amd64 := amd64
arch.arm := armhf
arch.arm64 := aarch64
arch.386 := 386

define gocross
	$(if $(findstring [$(1)/$(2)],$(VALID_OS_ARCH)), \
	echo "$(WARN_COLOR)$(MSG_PREFIX) Building binary for [$(1)/$(2)]$(MSG_SUFFIX)$(NO_COLOR)"; \
	echo "$(INFO_COLOR) GOOS=$(1) GOARCH=$(2) $(GO_ENV_FLAGS)\n $(GO) build -o $(BINARY_PATH)/$(BINARY_PREFIX)-${os.$(1)}-${arch.$(2)}$(call extension,$(GOOS))\n $(BUILD_FLAGS) . $(NO_COLOR)\n"; \
	GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 $(GO_ENV_FLAGS) \
		$(GO) build \
		-o $(BINARY_PATH)/$(BINARY_PREFIX)-${os.$(1)}-${arch.$(2)}$(call extension,$(GOOS)) \
		$(BUILD_FLAGS) .;)
endef

define buildTargets
	$(foreach GO_OS, $(3), $(foreach GO_ARCH, $(4), $(call $(1), $(2)$(GO_OS), -$(GO_ARCH))))
endef

build-x: $(shell find . -type f -name '*.go') ## to build for cross platforms.
	@$(foreach GO_OS, $(TARGET_PLATFORMS), $(foreach GO_ARCH, $(TARGET_ARCHS), $(call gocross,$(GO_OS),$(GO_ARCH))))

build-bin-for-%: $(shell find . -type f -name '*.go')
	@$(eval TARGET_PLATFORM=$(firstword $(subst -, , $*)))
	@$(eval TARGET_ARCH=$(or $(word 2,$(subst -, , $*)),$(value 2)))
	@$(call gocross,$(TARGET_PLATFORM),$(TARGET_ARCH))

build-bin: ## to build out a binary.
	$(if $(filter $(CROSS_BUILD), true), \
	    @$(MAKE) build-x, \
	    @$(MAKE) $(call buildTargets, addprefix, build-bin-for-, $(GO_OS), $(GO_ARCH)))

ifneq ($(BUILD_IN_CONTAINER), true)

build: build-bin ## to install dependencies and build out a binary.

else

build: deploy

endif

clean-bin: ## to clean generated binaries only.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Cleaning up binaries$(MSG_SUFFIX)$(NO_COLOR)"
	@for GO_OS in $(TARGET_PLATFORMS); do \
		for GO_ARCH in $(TARGET_ARCHS); do \
			TARGET_BINARY="$(BINARY_PATH)/$(BINARY_PREFIX)-$$GO_OS-$$GO_ARCH"; \
			if [ -f $$TARGET_BINARY ] ; then rm -Rf $$TARGET_BINARY ; fi \
		done; \
	done
