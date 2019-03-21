#!/usr/bin/env bash

set -eoux

CURRENT_DIR=$(dirname "$0")

. ${CURRENT_DIR}/functions.sh
. ${CURRENT_DIR}/common.sh

BASE_IMAGE=$BASE_IMAGE
PKG_NAMESPACE=$PKG_NAMESPACE
ARTIFACTS_PATH=/.artifacts
BINARY_PATH="${ARTIFACTS_PATH}/bin"
BINARY_PREFIX=$(basename `pwd` 2> /dev/null)

DOCKER_FILE=".dockerfile-${ARCH}"

if [[ $BASE_IMAGE == '' ]]; then
    BASE_IMAGE='scratch'
fi

if [[ $ARCH == '' ]]; then
    ARCH='amd64'
fi

cat ${CURRENT_DIR}/../templates/docker/Dockerfile.template.dist > $DOCKER_FILE

replaceInFile '{{BASE_IMAGE}}' $BASE_IMAGE $DOCKER_FILE
replaceInFile '{{ARCH}}' $ARCH $DOCKER_FILE
replaceInFile '{{PKG_NAMESPACE}}' $PKG_NAMESPACE $DOCKER_FILE
replaceInFile '{{ARTIFACTS_PATH}}' $ARTIFACTS_PATH $DOCKER_FILE
replaceInFile '{{BINARY_PATH}}' $BINARY_PATH $DOCKER_FILE
replaceInFile '{{BINARY_PREFIX}}' $BINARY_PREFIX $DOCKER_FILE
