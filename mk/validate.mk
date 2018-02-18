validate: vet lint format-check misspell ineffassign ## to validate the code.

# Simplified dead code detector. Used for skipping certain checks on unreachable code
# (for instance, shift checks on arch-specific code).
# https://golang.org/cmd/vet/
vet: ## to run detection on dead code.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Running vet$(MSG_SUFFIX)$(NO_COLOR)"
	@test -z "$$($(GO) vet $(GO_FLAGS) $(PKGS) 2>&1 | tee /dev/stderr)"

lint: ## to run linter against go files.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Running linter$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO_LINT) $(GO_FLAGS) $(PKGS)

format-check: ## to check if the go files are formatted correctly.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Checking code format$(MSG_SUFFIX)$(NO_COLOR)"
	@diff=$$($(GO_FMT) -d -s $(GO_FILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make format' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

format: ## to format all go files.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Formatting code$(MSG_SUFFIX)$(NO_COLOR)"
	@test -z "$$($(GO_FMT) -s -w $(GO_FLAGS) $(GO_FILES) 2>&1 | tee /dev/stderr)"

misspell: ## to run misspell.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Running misspell$(MSG_SUFFIX)$(NO_COLOR)"
	@test -z "$$(find . -type f | grep -v vendor/ | grep -v bin/ | grep -v .git/ | grep -v MAINTAINERS | xargs misspell | tee /dev/stderr)"

ineffassign: ## to run ineffassign.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Running ineffassign$(MSG_SUFFIX)$(NO_COLOR)"
	@test -z "$$(ineffassign . | grep -v vendor/ | grep -v ".pb.go:" | tee /dev/stderr)"
