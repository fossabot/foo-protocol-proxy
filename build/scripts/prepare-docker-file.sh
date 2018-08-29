#!/usr/bin/env bash

set -eoux

CURRENT_DIR=$(dirname "$0")
BASE_IMAGE=$BASE_IMAGE
ARCH=$ARCH
DOCKER_FILE=".dockerfile-${ARCH}"
SRC_NAMESPACE=$SRC_NAMESPACE
BINARY_PREFIX=$(basename `pwd` 2> /dev/null)

if [[ $BASE_IMAGE == '' ]]; then
    BASE_IMAGE='scratch'
fi

if [[ $ARCH == '' ]]; then
    ARCH='amd64'
fi

. ${CURRENT_DIR}/includes.sh

cat ${CURRENT_DIR}/../templates/docker/Dockerfile.template.dist > $DOCKER_FILE

replaceInFile '{{BASE_IMAGE}}' $BASE_IMAGE $DOCKER_FILE
replaceInFile '{{ARCH}}' $ARCH $DOCKER_FILE
replaceInFile '{{SRC_NAMESPACE}}' $SRC_NAMESPACE $DOCKER_FILE
replaceInFile '{{BINARY_PREFIX}}' $BINARY_PREFIX $DOCKER_FILE
