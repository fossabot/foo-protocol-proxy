COVERAGE_PACKAGES := app handlers persistence

define coverPackage
	$(GO) test -cover -parallel $(PARALLEL_TESTS) -timeout $(TEST_TIMEOUT) -covermode=$(COVERAGE_MODE) -coverprofile $(COVERAGE_PATH)/$(1).part $(GO_FLAGS) ./$(1);
endef

# Goveralls binary.
GOVERALLS_BIN := $(GOPATH)/bin/goveralls
GOVERALLS := $(shell [ -x $(GOVERALLS_BIN) ] && echo $(GOVERALLS_BIN) || echo '')

cover: $(COVERAGE_PROFILE) ## to run test with coverage and report that out to profile.

coverage-html: $(COVERAGE_HTML) ## to export coverage results to html format"$(COVERAGE_PATH)/index.html".
	@open "$(COVERAGE_HTML)"

# Send the results to coveralls.
coverage-send: $(COVERAGE_PROFILE)
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Sending coverage$(MSG_SUFFIX)$(NO_COLOR)"
	@$(if $(GOVERALLS), , $(error Please run make get-deps))
	@$(GOVERALLS) -service travis-ci -coverprofile="$(COVERAGE_PROFILE)" -repotoken $(COVERALLS_TOKEN)

coverage-serve: $(COVERAGE_HTML) ## to serve coverage results over http - useful only if building remote/headless
	@cd "$(COVERAGE_PATH)" && python -m SimpleHTTPServer 8888

# Clean up coverage output.
clean-coverage: ## to clean coverage generated files.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Cleaning up coverage generated files$(MSG_SUFFIX)$(NO_COLOR)"
	@rm -Rf "$(COVERAGE_PATH)"

$(COVERAGE_PROFILE):
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Coverage check$(MSG_SUFFIX)$(NO_COLOR)"
	@if [ ! -d $(COVERAGE_PATH) ] ; then mkdir -p $(COVERAGE_PATH) ; fi
	@$(foreach package, $(COVERAGE_PACKAGES), $(call coverPackage,$(package)))

	@echo "mode: $(COVERAGE_MODE)" > $(COVERAGE_PROFILE)
	@grep -h -v "mode: $(COVERAGE_MODE)" $(COVERAGE_PATH)/*.part >> "$(COVERAGE_PROFILE)"

$(COVERAGE_HTML): $(COVERAGE_PROFILE)
	@echo "$(WARN_COLOR)$(MSG_PREFIX) Coverage HTML export$(MSG_SUFFIX)$(NO_COLOR)"
	@if [ ! -d $(COVERAGE_PATH) ] ; then $(MAKE) $(COVERAGE_PATH) ; fi
	@$(GO) tool cover  -html="$(COVERAGE_PROFILE)" -o "$(COVERAGE_HTML)" $(GO_FLAGS)
