test: unit unit-short race bench integration ## to setup the external used tools.

unit: ## to run long unit tests.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Unit tests$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -cover -parallel $(PARALLEL_TESTS) -timeout=$(TEST_TIMEOUT) -tags $(GO_TAGS) $(GO_FLAGS) ./...

# Quick test. You can bypass long tests using: `if testing.Short() { t.Skip("Skipping in short mode.") }`.
unit-short: ## to run short unit tests.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Unit tests (short)$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -test.short -cover -parallel $(PARALLEL_TESTS) -timeout=$(TEST_TIMEOUT) -tags $(GO_TAGS) $(GO_FLAGS) ./...

# Runs long tests also, plus race detection.
race: ## to run long unit tests with race conditions coverage.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Unit tests with race cover$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -race -cpu=1,2,4 -parallel $(PARALLEL_TESTS) -timeout $(TEST_TIMEOUT) -tags $(GO_TAGS) $(GO_FLAGS) ./...

bench: ## to run benchmark tests.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Benchmarking tests$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -run NONE -bench . -benchmem -parallel $(PARALLEL_TESTS) -timeout $(TEST_TIMEOUT) -tags bench $(GO_FLAGS) $(PKGS)

integration: build ## to run integration tests.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Integration tests$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -cover -parallel $(PARALLEL_TESTS) -timeout $(TEST_TIMEOUT) -tags integration $(GO_FLAGS) ./...
