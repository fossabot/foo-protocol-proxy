validate: vet lint format-check ## to validate the code.

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
