DOCKER := $(shell docker info > /dev/null 2>&1 || $(SUDO) 2>&1) docker
DOCKER_LOGO := "🐳"

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
	@echo "$(INFO_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Preparing file$(MSG_SFX)$(NO_CLR)"
	@$(SUDO) BASE_IMAGE=${BASE_IMAGE}  \
	    ARCH=${ARCH}                   \
	    PKG_NAMESPACE=${PKG_NAMESPACE} \
	    bash $(DOCKER_FILE_SCRIPT_PATH) $(args) 2>&1

.env:
	@echo "$(INFO_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Preparing .env file$(MSG_SFX)$(NO_CLR)"
	@$(SUDO) ARCH=${ARCH} bash $(DOCKER_ENV_FILE_SCRIPT_PATH) $(args) 2>&1

build-dirs:
	@echo "$(WARN_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Building mapping directories for Go$(MSG_SFX)$(NO_CLR)"
	@mkdir -p ${GO_GENERATED_DIR}/bin ${GO_GENERATED_DIR}/pkg ${GO_GENERATED_DIR}/src/$(PKG_NAMESPACE) 2>&1

clean-containers: ## 🐳 to clean inactive containers data.
	@echo "$(WARN_CLR)$(MSG_PRFX) ${DOCKER_LOGO} 🧹 Cleaning up containers data$(MSG_SFX)$(NO_CLR)"
	@$(DOCKER) container prune -f 2>&1
	@$(eval EXITED_CONTAINERS := $(shell $(DOCKER) ps -aqf status=exited -f status=dead 2>&1))
	@test -n "${EXITED_CONTAINERS}" && $(DOCKER) rm ${EXITED_CONTAINERS} || true 2>&1

clean-files: ## 🐳 to clean deployment generated files.
	@echo "$(WARN_CLR)$(MSG_PRFX) ${DOCKER_LOGO} 🧹 Cleaning up generated files and directories$(MSG_SFX)$(NO_CLR)"
	@rm -rf .docker-compose-*.yaml .dockerfile* .env ${GO_GENERATED_DIR} 2>&1

clean-images: ## 🐳 to clean inactive images data.
	@echo "$(WARN_CLR)$(MSG_PRFX) ${DOCKER_LOGO} 🧹 Cleaning up images data$(MSG_SFX)$(NO_CLR)"
	@$(DOCKER) image prune -f 2>&1
	@$(eval DANGLING_IMAGES:= $(shell $(DOCKER) images -aqf dangling=true 2>&1))
	@test -n "${DANGLING_IMAGES}" && $(DOCKER) rmi ${DANGLING_IMAGES} || true 2>&1

clean-volumes: ## 🐳 to clean inactive containers volumes.
	@echo "$(WARN_CLR)$(MSG_PRFX) ${DOCKER_LOGO} 🧹 Cleaning up containers volumes$(MSG_SFX)$(NO_CLR)"
	@$(DOCKER) volume prune -f 2>&1
	@$(DANGLING_VOLUMES := $(shell $(DOCKER) volume ls -qf dangling=true 2>&1))
	@test -n "${DANGLING_VOLUMES}" && $(DOCKER) volume rm ${DANGLING_VOLUMES} || true 2>&1

deploy: build-dirs docker-prepare update-pkg-version ## 🐳 to deploy a docker container.
	@echo "$(INFO_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Deploying container$(MSG_SFX)$(NO_CLR)"
	@$(SUDO) ARCH=${ARCH} bash $(DEPLOY_SCRIPT_PATH) $(args) 2>&1

docker-clean: clean-images clean-containers clean-volumes ## 🐳 to clean inactive Docker data.

docker-exec: ## 🐳 to execute command inside the docker container.
	@echo "$(INFO_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Executing command inside the container$(MSG_SFX)$(NO_CLR)"
	@$(DOCKER) exec -it $(CONTAINER_NAME) $(CMD) 2>&1

docker-health: ## 🐳 to get the health state docker container.
	@echo "$(INFO_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Getting health state of the container$(MSG_SFX)$(NO_CLR)"
	@$(DOCKER) inspect --format='{{json .State.Health}}' $(CONTAINER_NAME) 2>&1

docker-kill: ## 🐳 to send kill signal to the main process at the docker container.
	@echo "$(WARN_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Sending kill($(args)) signal to main process$(MSG_SFX)$(NO_CLR)"
	@$(MAKE) docker-exec CMD="pkill $(args) $(BINARY_PREFIX)" > /dev/null 2>&1

docker-logs: ## 🐳 to get logs from the docker container.
	@echo "$(INFO_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Getting logs of the container$(MSG_SFX)$(NO_CLR)"
	@$(DOCKER) logs -f $(CONTAINER_NAME) 2>&1

docker-prepare: ## 🐳 prepare docker files from the templates.
	@echo "$(WARN_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Preparing docker files$(MSG_SFX)$(NO_CLR)"
	@$(SUDO) BASE_IMAGE=${BASE_IMAGE}              \
	    REGISTRY=${REGISTRY}                       \
	    IMAGE_NAME=${IMAGE_NAME}                   \
	    IMAGE_TAG=${DOCKER_IMAGE_TAG}              \
	    REGISTRY_REPO=${REGISTRY_REPO}             \
	    ARCH=${ARCH}                               \
	    SERVICE_NAME=${SERVICE_NAME}               \
	    SERVICE_DESCRIPTION=${SERVICE_DESCRIPTION} \
	    PKG_NAMESPACE=${PKG_NAMESPACE}             \
	    bash $(PREPARE_SCRIPT_PATH) $(args) 2>&1

docker-shell: .env build-dirs ## 🐳 run shell command inside the docker container.
	@echo "$(WARN_CLR)$(MSG_PRFX) ${DOCKER_LOGO} 🏃 running shell inside the container$(MSG_SFX)$(NO_CLR)"
	@$(DOCKER) run                                     \
	    -ti                                            \
	    --rm                                           \
	    -u $$(id -u):$$(id -g)                         \
	    -v "$$(pwd)/${GO_GENERATED_DIR}:/go"           \
	    -v "$$(pwd):/go/src/$(PKG_NAMESPACE)"          \
	    -w "/go/src/$(PKG_NAMESPACE)"                  \
	    -v "$$(pwd)/${GO_GENERATED_DIR}/bin:/go/bin"   \
	    -v "$$(pwd)/${GO_GENERATED_DIR}/cache:/.cache" \
	    -w /go/src/$(PKG_NAMESPACE)                    \
	    --env-file .env                                \
	    $(BUILD_IMAGE)                                 \
	    /bin/sh $(CMD) 2>&1

image: .dockerfile-$(ARCH) ## 🐳 to build a docker image.
	@echo "$(INFO_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Creating image \"$(BUILD_IMAGE)\"$(MSG_SFX)$(NO_CLR)"
	@$(DOCKER) build          \
	    ${DOCKER_BUILD_FLAGS} \
	    -t $(BUILD_IMAGE)     \
	    -f $(DOCKER_FILE)     \
	    $(args) . 2>&1

list-containers: ## 🐳 to list all containers.
	@echo "$(INFO_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Listing containers$(MSG_SFX)$(NO_CLR)"
	@$(DOCKER) ps -a --format "{{.ID}} $(MSG_PRFX) {{.Names}}" 2>&1

publish: ## 🐳 to publish the docker image to dockerhub repository.
	@echo "$(WARN_CLR)$(MSG_PRFX) ${DOCKER_LOGO} Pushing image to \"$(BUILD_IMAGE)$\"(MSG_SFX)$(NO_CLR)"
	@$(DOCKER) push $(BUILD_IMAGE) 2>&1
