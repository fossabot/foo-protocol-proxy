GO_LINT ?= golint
GO_FMT ?= gofmt

validate: ineffassign format-check lint misspell vet ## to validate the code.

ineffassign: ## to run ineffassign.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Running ineffassign$(MSG_SUFFIX)$(NO_COLOR)"
	@test -z "$$(ineffassign . | grep -v vendor/ | grep -v ".pb.go:" | tee /dev/stderr)"

format-check: ## to check if the Go files are formatted correctly.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Checking code format$(MSG_SUFFIX)$(NO_COLOR)"
	@diff=$$($(GO_FMT) -d -s $(GO_FILES)); \
	if [ -n "$$diff" ]; then                                    \
		echo "Please run 'make format' and commit the result:"; \
		echo "$${diff}";                                        \
		exit 1;                                                 \
	fi;

lint: ## to run linter against Go files.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Running linter$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO_LINT) $(GO_FLAGS) $(SRC_PKGS)

misspell: ## to run misspell.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Running misspell$(MSG_SUFFIX)$(NO_COLOR)"
	@test -z "$$(find . -type f | grep -v vendor/ | grep -v bin/ | grep -v .git/ | grep -v MAINTAINERS | xargs misspell | tee /dev/stderr)"

# Simplified dead code detector. Used for skipping certain checks on unreachable code
# (for instance, shift checks on arch-specific code).
# https://golang.org/cmd/vet/
vet: ## to run detection on dead code.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Running vet$(MSG_SUFFIX)$(NO_COLOR)"
	@test -z "$$($(GO) vet $(GO_FLAGS) $(SRC_PKGS) 2>&1 | tee /dev/stderr)"
