DOCKER := $(shell docker info > /dev/null 2>&1 || $(SUDO) 2>&1) docker

image: ## to build a docker image.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) 🐳 Creating Docker Image$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) build ${DOCKER_BUILD_FLAGS} -t $(REGISTRY_REPO):$(DOCKER_TAG) -f $(DOCKER_FILE) $(args) 2>&1

deploy: ## to deploy a docker container.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) 🐳 Deploying Docker Container$(MSG_SUFFIX)$(NO_COLOR)"
	@$(SUDO) bash ./deploy.sh $(args) 2>&1

publish: ## to publish the docker image to dockerhub repository.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) 🐳 Pushing Docker Image to $(REGISTRY_REPO):$(DOCKER_TAG)$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) push $(REGISTRY_REPO):$(DOCKER_TAG) 2>&1

docker-kill: ## to send kill signal to the main process at the docker container.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) 🐳 Sending kill signal to main Docker process$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) exec -it $(BINARY_PREFIX)-${DOCKER_TAG} pkill $(args) $(BINARY_PREFIX) > /dev/null 2>&1

docker-logs: ## to get logs from the docker container.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) 🐳 Getting logs of the Docker container$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) logs -f $(BINARY_PREFIX)-${DOCKER_TAG} 2>&1
