#!/usr/bin/env bash

set -eoux

CURRENT_DIR=$(dirname "$0")

. ${CURRENT_DIR}/functions.sh
. ${CURRENT_DIR}/common.sh

docker-compose --verbose -f $DOCKER_COMPOSE_FILE build --force-rm --pull --no-cache
docker-compose --verbose -f $DOCKER_COMPOSE_FILE up -d
