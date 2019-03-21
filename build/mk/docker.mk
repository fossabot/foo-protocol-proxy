DOCKER := $(shell docker info > /dev/null 2>&1 || $(SUDO) 2>&1) docker
DOCKER_LOGO := "🐳"

SRC_NAMESPACE := $(shell $(GO) list $(COMMAND_DIR) 2> /dev/null | sed -e "s|$(subst .,,$(COMMAND_DIR))||g")

images.amd64 = alpine
images.arm64 = alpine
images.arm = alpine
images.ppc64le = alpine

BASE_IMAGE ?= ${images.$(ARCH)}

## 🐳 Target build image name.
BUILD_IMAGE ?= $(REGISTRY_REPO):$(DOCKER_IMAGE_TAG)
## 🐳 Container name that the action will be performed on.
CONTAINER_NAME ?= $(BINARY_PREFIX).${DOCKER_IMAGE_TAG}

all-docker: build-dirs clean-containers clean-images clean-volumes deploy docker-clean docker-health image list-containers

.dockerfile-$(ARCH): .env
	@echo "$(INFO_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Preparing Docker file$(MSG_SUFFIX)$(NO_COLOR)"
	@$(SUDO) BASE_IMAGE=${BASE_IMAGE} ARCH=${ARCH} SRC_NAMESPACE=${SRC_NAMESPACE} \
	bash $(DOCKER_FILE_SCRIPT_PATH) $(args) 2>&1

.env:
	@echo "$(INFO_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Preparing Docker .env file$(MSG_SUFFIX)$(NO_COLOR)"
	@$(SUDO) bash $(DOCKER_ENV_FILE_SCRIPT_PATH) $(args) 2>&1

build-dirs:
	@echo "$(WARN_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Building mapping directories for Go$(MSG_SUFFIX)$(NO_COLOR)"
	@mkdir -p .go .go/src/$(PKG_BASE) .go/bin .go/pkg

clean-containers: ## 🐳 to clean inactive containers data.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Cleaning up Docker containers data$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) container prune -f 2>&1
	@$(eval EXITED_CONTAINERS := $(shell $(DOCKER) ps -aqf status=exited -f status=dead 2>&1))
	@test -n "${EXITED_CONTAINERS}" && $(DOCKER) rm ${EXITED_CONTAINERS} || true

clean-files: ## 🐳 to clean deployment generated files.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Cleaning up Docker generated files$(MSG_SUFFIX)$(NO_COLOR)"
	@rm -rf .docker-compose-*.yaml .dockerfile* .env .go 2>&1

clean-images: ## 🐳 to clean inactive images data.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Cleaning up Docker images data$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) image prune -f 2>&1
	@$(eval DANGLING_IMAGES:= $(shell $(DOCKER) images -aqf dangling=true 2>&1))
	@test -n "${DANGLING_IMAGES}" && $(DOCKER) rmi ${DANGLING_IMAGES} || true

clean-volumes: ## 🐳 to clean inactive containers volumes.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Cleaning up Docker containers volumes$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) volume prune -f 2>&1
	@$(DANGLING_VOLUMES := $(shell $(DOCKER) volume ls -qf dangling=true 2>&1))
	@test -n "${DANGLING_VOLUMES}" && $(DOCKER) volume rm ${DANGLING_VOLUMES} || true

deploy: build-dirs ## 🐳 to deploy a docker container.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Deploying Docker container$(MSG_SUFFIX)$(NO_COLOR)"
	@$(SUDO) BASE_IMAGE=${BASE_IMAGE} REGISTRY=${REGISTRY} IMAGE_NAME=${IMAGE_NAME} IMAGE_TAG=${DOCKER_IMAGE_TAG} \
	 REGISTRY_REPO=${REGISTRY_REPO} ARCH=${ARCH} SERVICE_NAME=${SERVICE_NAME}                                     \
	 SERVICE_DESCRIPTION=${SERVICE_DESCRIPTION} SRC_NAMESPACE=${SRC_NAMESPACE}                                    \
	 bash $(DEPLOY_SCRIPT_PATH) $(args) 2>&1

docker-clean: clean-images clean-containers clean-volumes ## 🐳 to clean inactive Docker data.

docker-exec: ## 🐳 to execute command inside the docker container.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Executing command inside the Docker container$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) exec -it $(CONTAINER_NAME) $(CMD) 2>&1

docker-health: ## 🐳 to get the health state docker container.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Getting health state of the Docker container$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) inspect --format='{{json .State.Health}}' $(CONTAINER_NAME) 2>&1

docker-logs: ## 🐳 to get logs from the docker container.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Getting logs of the Docker container$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) logs -f $(CONTAINER_NAME) 2>&1

docker-kill: ## 🐳 to send kill signal to the main process at the docker container.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Sending kill($(args)) signal to main Docker process$(MSG_SUFFIX)$(NO_COLOR)"
	@$(MAKE) docker-exec CMD="pkill $(args) $(BINARY_PREFIX)" > /dev/null 2>&1

docker-run: .env build-dirs
	@echo "$(WARN_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Building temp Docker container$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) run                       \
	    -ti                              \
	    --rm                             \
	    -u $$(id -u):$$(id -g)           \
	    -v "$$(pwd)/.go:/go"             \
	    -v "$$(pwd):/go/src/$(PKG_BASE)" \
	    -v "$$(pwd)/.bin:/go/bin"        \
	    -v "$$(pwd)/.go/cache:/.cache"   \
	    -w /go/src/$(PKG_BASE)           \
	    --env-file .env                  \
	    $(BUILD_IMAGE)                   \
	    /bin/sh $(CMD)

image: .dockerfile-$(ARCH) ## 🐳 to build a docker image.
	@echo "$(INFO_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Creating Docker image \"$(BUILD_IMAGE)\"$(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) build ${DOCKER_BUILD_FLAGS} -t $(BUILD_IMAGE) -f $(DOCKER_FILE) $(args) . 2>&1

list-containers: ## 🐳 to list all containers.
	@$(DOCKER) ps -a --format "{{.ID}} $(MSG_PREFIX) {{.Names}}" 2>&1

publish: ## 🐳 to publish the docker image to dockerhub repository.
	@echo "$(WARN_COLOR)$(MSG_PREFIX) ${DOCKER_LOGO} Pushing Docker image to \"$(BUILD_IMAGE)$\"(MSG_SUFFIX)$(NO_COLOR)"
	@$(DOCKER) push $(BUILD_IMAGE) 2>&1
