COVERAGE_PACKAGES := app handlers persistence

define coverPackage
	$(GO) test -cover -parallel $(PARALLEL_TESTS) -timeout $(TEST_TIMEOUT) -covermode=$(COVERAGE_MODE) -coverprofile $(COVERAGE_PATH)/$(1).part $(GO_FLAGS) ./$(1) 2>&1;
endef

# Goveralls binary.
GOVERALLS_BIN := $(GOPATH)/bin/goveralls
GOVERALLS := $(shell [ -x $(GOVERALLS_BIN) ] && echo $(GOVERALLS_BIN) || echo '' 2> /dev/null)

all-coverage: cover coverage-browse

$(COVERAGE_HTML): $(COVERAGE_PROFILE)
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Coverage HTML export$(MSG_SUFFIX)$(NO_COLOR)"
	@if [ ! -d $(COVERAGE_PATH) ] ; then $(MAKE) $(COVERAGE_PATH) 2>&1 ; fi
	@$(GO) tool cover -html="$(COVERAGE_PROFILE)" -o "$(COVERAGE_HTML)" $(GO_FLAGS) 2>&1

$(COVERAGE_PROFILE):
	@if [ ! -d $(COVERAGE_PATH) ] ; then mkdir -p $(COVERAGE_PATH) 2>&1 ; fi
	@$(foreach package, $(COVERAGE_PACKAGES), $(call coverPackage,$(package)))

	@echo "mode: $(COVERAGE_MODE)" > $(COVERAGE_PROFILE) 2>&1
	# tail -q -n +2 $(COVERAGE_PATH)/*.part
	@grep -h -v "mode: $(COVERAGE_MODE)" $(COVERAGE_PATH)/*.part >> "$(COVERAGE_PROFILE)" 2>&1

# Clean up coverage output.
clean-coverage: ## to clean coverage generated files.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Cleaning up coverage generated files$(MSG_SUFFIX)$(NO_COLOR)"
	@rm -rf "$(COVERAGE_PATH)" 2>&1

cover: $(COVERAGE_PROFILE) ## to run test with coverage and report that out to profile.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Coverage check$(MSG_SUFFIX)$(NO_COLOR)"

coverage-browse: $(COVERAGE_HTML) ## to export coverage results to html format"$(COVERAGE_PATH)/index.html".
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Opening browser$(MSG_SUFFIX)$(NO_COLOR)"
	@open "$(COVERAGE_HTML)" 2>&1

# Send the results to coveralls.
coverage-send: $(COVERAGE_PROFILE)
	@echo "$(INFO_COLOR)$(MSG_PREFIX) Sending coverage$(MSG_SUFFIX)$(NO_COLOR)"
	@$(if $(GOVERALLS), , $(error Please run make get-deps))
	@$(GOVERALLS) -service travis-ci -coverprofile="$(COVERAGE_PROFILE)" -repotoken $(COVERALLS_TOKEN) 2>&1

coverage-serve: $(COVERAGE_HTML) ## to serve coverage results over http - useful only if building remote/headless.
	@cd "$(COVERAGE_PATH)" && python -m SimpleHTTPServer 8888 2>&1
