test: bench integration profile race unit unit-short ## to setup the external used tools.

$(BENCH_TESTS_PATH):
	@if [ ! -d $(BENCH_TESTS_PATH) ] ; then mkdir -p $(BENCH_TESTS_PATH) 2>&1 ; fi

bench: ## to run benchmark tests.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Benchmarking tests$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -count=$(BENCH_TESTS_COUNT) -run=NONE -bench . -benchmem -parallel $(PARALLEL_TESTS) -timeout $(TEST_TIMEOUT) -tags bench $(GO_FLAGS) $(SRC_PKGS) 2>&1

# Clean up bench tests output.
clean-bench: ## to clean coverage generated files.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Cleaning up bench tests generated files$(MSG_SUFFIX)$(NO_COLOR)"
	@rm -rf "$(BENCH_TESTS_PATH)" 2>&1

# Clean up tests output.
clean-tests: clean-bench ## to clean coverage generated files.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Cleaning up tests generated files$(MSG_SUFFIX)$(NO_COLOR)"
	@rm -rf "$(TESTS_PATH)" 2>&1

integration: build ## to run integration tests.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Integration tests$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -cover -parallel $(PARALLEL_TESTS) -timeout $(TEST_TIMEOUT) -tags integration $(GO_FLAGS) ./... 2>&1

profile: $(BENCH_TESTS_PATH) ## to get bench mark profiles.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Bench tests check$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -count=$(BENCH_TESTS_COUNT) -run=NONE -bench . -benchmem -parallel $(PARALLEL_TESTS) -timeout $(TEST_TIMEOUT) -tags bench -o $(BENCH_TESTS_PATH)/test.bin -cpuprofile $(BENCH_TESTS_PATH)/cpu.out -memprofile $(BENCH_TESTS_PATH)/mem.out $(GO_FLAGS) . 2>&1
	@$(GO) tool pprof --svg $(BENCH_TESTS_PATH)/test.bin $(BENCH_TESTS_PATH)/mem.out > $(BENCH_TESTS_PATH)/mem.svg
	@$(GO) tool pprof --svg $(BENCH_TESTS_PATH)/test.bin $(BENCH_TESTS_PATH)/cpu.out > $(BENCH_TESTS_PATH)/cpu.svg

# Runs long tests also, plus race detection.
race: ## to run long unit tests with race conditions coverage.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Unit tests with race cover$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -race -cpu=1,2,4 -parallel $(PARALLEL_TESTS) -timeout $(TEST_TIMEOUT) -tags $(GO_TAGS) $(GO_FLAGS) ./... 2>&1

unit: ## to run long unit tests.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Unit tests$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -cover -parallel $(PARALLEL_TESTS) -timeout=$(TEST_TIMEOUT) -tags $(GO_TAGS) $(GO_FLAGS) ./... 2>&1

# Quick test. You can bypass long tests using: `if testing.Short() { t.Skip("Skipping in short mode.") }`.
unit-short: ## to run short unit tests.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Unit tests (short)$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) test -test.short -cover -parallel $(PARALLEL_TESTS) -timeout=$(TEST_TIMEOUT) -tags $(GO_TAGS) $(GO_FLAGS) ./... 2>&1
