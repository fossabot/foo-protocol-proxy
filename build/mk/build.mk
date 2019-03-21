GO_OS ?= $(shell $(GO) env GOOS 2> /dev/null)
GO_ARCH ?= $(shell $(GO) env GOARCH 2> /dev/null)

# Disable usage of CGO.
CGO_ENABLED := 0

ifeq ($(ARCH),)
    ARCH := $(GO_ARCH)
endif

ifeq ($(OS),)
    OS := $(GO_OS)
endif

# Valid OS and Architecture target combinations. Check https://golang.org/doc/install/source#environment
# Or run `go tool dist list -json | jq`
VALID_OS_ARCH := "[darwin/amd64][linux/amd64][linux/arm64][linux/arm][linux/ppc64le][openbsd/amd64][windows/amd64][windows/386]"

os.darwin := Darwin
os.linux := Linux
os.openbsd := OpenBSD
os.windows := Windows

arch.amd64 := x86_64
arch.386 := i386
arch.arm64 := aarch64
arch.arm := armhf
arch.ppc64le := ppc64le
arch.s390x := s390x

BINARY_NAME := $(BINARY_PREFIX)-${os.$(OS)}-${arch.$(ARCH)}
TARGET_BINARY := $(BINARY_PATH)/$(BINARY_NAME)

# Package main path.
PKG_BASE ?= $(shell $(GO) list -e ./ 2> /dev/null)

# Get the current local branch name from git (if we can, this may be blank).
GIT_BRANCH := $(shell git symbolic-ref --short HEAD 2> /dev/null || git rev-parse --abbrev-ref HEAD 2> /dev/null)
GIT_COMMIT := $(shell git rev-parse HEAD 2> /dev/null)
# Get the git commit.
GIT_DIRTY := $(shell test -n "`git status --porcelain --untracked-files=no 2> /dev/null`" && echo "+CHANGES" || true 2> /dev/null)

# Build Flags
# The default version that's chosen when pushing the images. Can/should be overridden.
BUILD_VERSION ?= $(shell git describe --abbrev=8 --dirty='-Changes' 2> /dev/null | cut -d "v" -f 2 2> /dev/null)
BUILD_HASH ?= git-$(shell git rev-parse --short=18 HEAD 2> /dev/null)
BUILD_TIME ?= $(shell date +%FT%T%z 2> /dev/null)

# If we don't set the build version it defaults to dev.
ifeq ($(BUILD_VERSION),)
	BUILD_VERSION := $(shell cat $(CURDIR)/.version 2> /dev/null || echo dev)
endif

BUILD_ENV ?= $(BUILD_ENV:)

GO_BUILD_FLAGS ?= -i -a -installsuffix cgo

EXTLD_FLAGS ?=

# Check if we are not building for darwin, and honoring static.
IS_DARWIN_HOST ?= $(shell echo $(GO_OS) | egrep -i -c "darwin" 2> /dev/null)
IS_STATIC ?= $(shell echo $(STATIC) | egrep -i -c "true" 2>&1)

# Below, we are building a boolean circuit that says "$(IS_DARWIN_HOST) && $(IS_STATIC)"
ifeq ($(shell echo $$(( $(IS_DARWIN_HOST) * $(IS_STATIC) )) 2> /dev/null), 0)
# The flags we are passing to go build. -extldflags -static for making a static binary,
# or -linkmode external for linking external C libraries into the binary.
    override EXTLD_FLAGS +=-lm -static -lstdc++ -lpthread -static-libstdc++
endif

# -X version.BuildHash for telling the Go binary which build hash is used in this version,
# -X version.BuildTime for telling the Go binary the build time,
# -X version.GitBranch for telling the Go binary the git branch used,
# -X version.GitCommit for telling the Go binary the git commit used,
# -X main.version for telling the Go binary which version it is.
GO_LINKER_FLAGS ?=-s                                   \
        -v                                             \
        -w                                             \
        -X ${PKG_BASE}/version.BuildHash=$(BUILD_HASH) \
        -X ${PKG_BASE}/core.BuildTime=$(BUILD_TIME)    \
        -X ${PKG_BASE}/core.GitBranch=$(GIT_BRANCH)    \
        -X ${PKG_BASE}/core.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)

ifdef BUILD_VERSION
	GO_LINKER_FLAGS += -X main.version=$(BUILD_VERSION)
	DOCKER_IMAGE_TAG = $(BUILD_VERSION)
endif

ifdef EXTLD_FLAGS
    GO_LINKER_FLAGS	+= -extldflags "$(EXTLD_FLAGS)"
endif

GO_GC_FLAGS :=-trimpath=$(CURDIR)

# Honor debug
ifeq ($(DEBUG), true)
	# Disable function inlining and variable registration.
	GO_GC_FLAGS +=-N -l
endif

GO_ASM_FLAGS :=-trimpath=$(CURDIR)

# netgo for enforcing the native Go DNS resolver.
GO_TAGS ?= netgo

GO_ENV_FLAGS ?= $(GO_ENV_FLAGS:)
GO_ENV_FLAGS += $(BUILD_ENV)

extension = $(patsubst windows, .exe, $(filter windows, $(1)))

define goCross
	$(if $(findstring [$(1)/$(2)],$(VALID_OS_ARCH)),                                             \
	printf "$(OK_COLOR)$(MSG_PREFIX) Building binary for [$(1)/$(2)]$(MSG_SUFFIX)$(NO_COLOR)\n"; \
	printf "$(INFO_COLOR)\
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(1) GOARCH=$(2) $(GO_ENV_FLAGS)\n\
	    $(GO) build -o $(BINARY_PATH)/$(BINARY_PREFIX)-${os.$(1)}-${arch.$(2)}$(call extension,$(GO_OS))\n\
	    $(GO_BUILD_FLAGS)\n\
	    -ldflags '$(shell echo $(GO_LINKER_FLAGS) | sed -e 's|extldflags $(EXTLD_FLAGS)|extldflags \\"$(EXTLD_FLAGS)\\"|g' 2> /dev/null)'\n\
	    -gcflags=\"$(GO_GC_FLAGS)\"\n\
	    -asmflags=\"$(GO_ASM_FLAGS)\"\n\
	    -tags $(GO_TAGS)\n\
	    $(GO_FLAGS) $(COMMAND_DIR)\
	    $(NO_COLOR)\n";                                                                         \
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(1) GOARCH=$(2) $(GO_ENV_FLAGS)                            \
		$(GO) build                                                                             \
		-o $(BINARY_PATH)/$(BINARY_PREFIX)-${os.$(1)}-${arch.$(2)}$(call extension,$(GO_OS))    \
		$(GO_BUILD_FLAGS)                                                                       \
		-ldflags '$(GO_LINKER_FLAGS)'                                                           \
		-gcflags="$(GO_GC_FLAGS)"                                                               \
		-asmflags="$(GO_ASM_FLAGS)"                                                             \
		-tags $(GO_TAGS)                                                                        \
		$(GO_FLAGS) $(COMMAND_DIR);,                                                            \
		printf "$(ERROR_COLOR)Not defined build target \"[$(1)/$(2)]\"$(NO_COLOR)\n";           \
		printf "$(INFO_COLOR)Defined build tragets are: $(VALID_OS_ARCH).$(NO_COLOR)\n"         \
	)
endef

define buildTargets
	$(foreach GO_OS, $(3), $(foreach GO_ARCH, $(4), $(call $(1), $(2)$(GO_OS), -$(GO_ARCH))))
endef

define getDependency
	$(GO) get -u -v $(GO_FLAGS) $(1) 2>&1;
endef

define replaceInFile
    $(if $(findstring $(IS_DARWIN_HOST),1),  \
        sed -i '' "s|$(1)|$(2)|g" $(3) 2>&1, \
        sed -i -e "s|$(1)|$(2)|g" $(3) 2>&1)
endef

all-build: build-bin build-version build-x clean-bin clean-version get-deps go-generate go-install install update-pkg-version uninstall version

ifneq ($(BUILD_IN_CONTAINER), true)

build: build-bin ## to install dependencies and build out a binary.

else

build: deploy

endif

build-bin: ## to build out a binary.
	$(if $(filter $(CROSS_BUILD), true), \
	    @$(MAKE) build-x,                \
	    @$(MAKE) $(call buildTargets, addprefix, build-bin-for-, $(OS), $(ARCH)))

build-bin-for-%: $(shell find . -type f -name '*.go')
	@$(eval TARGET_PLATFORM=$(firstword $(subst -, , $*)))
	@$(eval TARGET_ARCH=$(or $(word 2,$(subst -, , $*)),$(value 2)))
	@$(call goCross,$(TARGET_PLATFORM),$(TARGET_ARCH))

build-version:  ## to get the current build version.
	@echo $(BUILD_VERSION)

build-x: $(shell find . -type f -name '*.go') ## to build for cross platforms.
	@$(foreach GO_OS, $(TARGET_PLATFORMS), $(foreach GO_ARCH, $(TARGET_ARCHS), $(call goCross,$(GO_OS),$(GO_ARCH))))

clean-bin: ## to clean generated binaries only.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Cleaning up binaries$(MSG_SUFFIX)$(NO_COLOR)"
	@rm -rf $(BINARY_PATH) 2>&1

clean-version: ## to remove updated go package version.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Cleaning up updated package version$(MSG_SUFFIX)$(NO_COLOR)"
	@rm -rf  $(PKG_TEMPLATE_DIR)/version.go

get-deps: ## to get required dependencies.
	@printf "$(OK_COLOR)$(MSG_PREFIX) Installing required dependencies$(MSG_SUFFIX)$(NO_COLOR)\n"
	@$(GO) mod download $(GO_FLAGS) 2>&1
	@$(foreach dependency, $(DEPENDENCIES), $(call getDependency,$(dependency)))

go-generate: ## to generate Go related files.
	@echo "$(OK_COLOR)$(MSG_PREFIX) Generating files via Go generate$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) generate $(GO_FLAGS) $(SRC_PKGS) 2>&1

go-install: update-pkg-version ## to install the Go related/dependent commands and packages.
	@printf "$(OK_COLOR)$(MSG_PREFIX) Installing Go related dependencies$(MSG_SUFFIX)$(NO_COLOR)\n"
	@$(GO) install -ldflags "-X $(PKG_BASE)/pkg/version.VERSION=${VERSION}" \
	 -installsuffix "static"                                                \
	 -v $(GO_FLAGS)                                                         \
	 -tags $(GO_TAGS)                                                       \
	 $(SRC_PKGS) 2>&1

install: ## to install the generated binary.
	@echo "$(OK_COLOR)$(MSG_PREFIX) Installing generated binary$(MSG_SUFFIX)$(NO_COLOR)"
	@if [ ! -f $(TARGET_BINARY) ] ; then $(MAKE) build; fi
	@cp $(TARGET_BINARY) $(INSTALLATION_BASE_PATH) 2>&1

kill: ## to send a kill signal to the running process of the binary.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Sending kill signal $(args)$(MSG_SUFFIX)$(NO_COLOR)"
	@pkill $(args) $(notdir $(TARGET_BINARY)) > /dev/null 2>&1

run: ## to run the generated binary, and build a new one if not existed.
	@echo "$(OK_COLOR)$(MSG_PREFIX) Running generated binary$(MSG_SUFFIX)$(NO_COLOR)"
	@if [ ! -f $(TARGET_BINARY) ] ; then $(MAKE) build; fi
	@$(TARGET_BINARY) $(args) 2>&1

uninstall: ## to uninstall generated binary.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Uninstalling generated binary$(MSG_SUFFIX)$(NO_COLOR)"
	@rm -rf $(INSTALLATION_BASE_PATH)/$(BINARY_NAME) 2>&1

update-pkg-version: ## to update package version.
	@printf "$(INFO_COLOR)$(MSG_PREFIX) Updating Go package version$(MSG_SUFFIX)$(NO_COLOR)\n"
    ifneq ($(wildcard $(PKG_TEMPLATE_DIR)/$(PKG_TEMPLATE)),)
		@cp $(PKG_TEMPLATE_DIR)/$(PKG_TEMPLATE) $(PKG_TEMPLATE_DIR)/version.go 2>&1
		@$(call replaceInFile,<VERSION>,$(VERSION),$(PKG_TEMPLATE_DIR)/version.go)
    endif

version:  ## to get the current version.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Current tagged version$(MSG_SUFFIX)$(NO_COLOR)"
	@echo "$(OK_COLOR)$(MSG_PREFIX) $(VERSION) $(NO_COLOR)"
