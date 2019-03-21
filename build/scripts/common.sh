#!/usr/bin/env bash

CURRENT_DIR=$(dirname "$0")

ARCH=$ARCH

# Environment variables.
FORWARDING_PORT=8001
LISTENING_PORT=8002
HTTP_PORT=8088
HEALTH_PORT=8081
RECOVERY_PATH='.data/recovery.json'

ENV_FILE_PATH='.env'
DOCKER_COMPOSE_FILE=".docker-compose-${ARCH}.yaml"
