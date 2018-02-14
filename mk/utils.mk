all: setup generate test coverage-html validate format nuke build-x install clean list help deploy ## to run all targets.

list: ## to list all targets.
	@$(MAKE) -rRpqn | awk -F':' '/^[a-z0-9][^$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)printf "$(DISCLAIMER_COLOR)%-30s$(NO_COLOR)\n", A[i]}' | sort -u 2>&1

help: ## to get help about the targets.
	@echo "$(OK_COLOR)$$FOO_PROTOCOL_PROXY$(NO_COLOR)"
	@echo "$(INFO_COLOR)Please use \`make <target>\`, Available options for <target> are:$(NO_COLOR)\n"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z._-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "  $(HELP_COLOR)%-28s$(NO_COLOR)  %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort -u 2>&1
	@echo "Useful variables:"
	@awk 'BEGIN { FS = ":=" } /^## /{x = substr($$0, 4); getline; if (NF >= 2) printf "  $(DISCLAIMER_COLOR)%-28s$(NO_COLOR) %s\n", $$1, x}' $(MAKEFILE_LIST) | sort -u 2>&1

setup: ## to setup the external used tools.
	@echo "$(OK_COLOR)$(MSG_PREFIX) Setting-up required components$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) install $(GO_FLAGS) -tags $(GO_TAGS) $(PKGS) 2>&1

get-deps: ## to get required dependencies.
	@echo "$(OK_COLOR)$(MSG_PREFIX) Installing required dependencies$(MSG_SUFFIX)$(NO_COLOR)"
	@$(foreach dependency, $(DEPENDENCIES), $(call getDependency,$(dependency)))

generate: ## to generate related files.
	@echo "$(OK_COLOR)$(MSG_PREFIX) Generating files via go generate$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) generate $(GO_FLAGS) $(PKGS) 2>&1

install: ## to install the generated binary.
	@echo "$(OK_COLOR)$(MSG_PREFIX) Installing generated binary$(MSG_SUFFIX)$(NO_COLOR)"
	@if [ ! -f $(TARGET_BINARY) ] ; then $(MAKE) build; fi
	@cp $(TARGET_BINARY) $(INSTALLATION_BASE_PATH) 2>&1

clean: clean-bin clean-coverage ## to clean up all generated files.

nuke: clean ## to do clean up and enforce removing the corresponding installed archive or binary.
	@echo "$(OK_COLOR)$(MSG_PREFIX) Cleaning go dependencies$(MSG_SUFFIX)$(NO_COLOR)"
	@$(GO) clean -i $(GO_FLAGS) ./... net 2>&1

run: ## to run the generated binary, and build a new one if not existed.
	@echo "$(OK_COLOR)$(MSG_PREFIX) Running generated binary$(MSG_SUFFIX)$(NO_COLOR)"
	@if [ ! -f $(TARGET_BINARY) ] ; then $(MAKE) build; fi
	@$(TARGET_BINARY) $(args) 2>&1

kill: ## to send a kill signal to the running process of the binary.
	@echo "$(OK_COLOR)$(MSG_PREFIX) Sending kill signal $(args)$(MSG_SUFFIX)$(NO_COLOR)"
	@pkill $(args) $(notdir $(TARGET_BINARY)) > /dev/null 2>&1
