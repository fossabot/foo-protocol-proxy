#!/usr/bin/env bash

set -eux

IMAGE_PREFIX="ahmedkamal"
REGISTRY_REPO="foo-protocol-proxy"
DOCKER_FILE="Dockerfile-proxy"
IMAGE_TAG=$(git describe --abbrev=0 | sed 's/^v\(.*\)$/\1/' 2> /dev/null)
FORWARDING_PORT=8001
LISTENING_PORT=8002
HTTP_PORT=8088
RECOVERY_PATH="data/recovery.json"

while getopts "f:l:h:r:p:t:" name
do
    case $name in
    f)  FORWARDING_PORT="$OPTARG";;
    l)  LISTENING_PORT="$OPTARG";;
    h)  HTTP_PORT="$OPTARG";;
    r)  RECOVERY_PATH="$OPTARG";;
    p)  IMAGE_PREFIX="$OPTARG";;
    t)  IMAGE_TAG="$OPTARG";;
    ?)  printf "Usage: %f: [-f forwarding port] [-l listening port] [-h http address] [-r recovery path] [-p image prefix] [-t image tag]\n" $0
          exit 2;;
    esac
done

if [[ $IMAGE_PREFIX == '' ]]; then
    IMAGE_PREFIX="ahmedkamal"
fi

if [[ $IMAGE_TAG == '' ]]; then
    IMAGE_TAG="dev"
fi

export IMAGE_PREFIX="$IMAGE_PREFIX"
export REGISTRY_REPO="$REGISTRY_REPO"
export DOCKER_FILE="$DOCKER_FILE"
export IMAGE_TAG="$IMAGE_TAG"
export FORWARDING_PORT="$FORWARDING_PORT"
export LISTENING_PORT="$LISTENING_PORT"
export HTTP_PORT="$HTTP_PORT"
export RECOVERY_PATH="$RECOVERY_PATH"

# Gets the current platform.
getPlatform() {
  echo $(uname -s | tr '[:upper:]' '[:lower:]') 2>&1
}

# Should an ip alias used based on the platform.
shouldUseIpAlias() {
  local platform=$1
  if [[ $platform == 'linux' ]]; then
     return 1
  elif [[ $platform == 'darwin' ]]; then
     return 0
  fi

  return 1
}

export HOST_OS=$(getPlatform)
if ( shouldUseIpAlias $HOST_OS )
then
  export HOST_IP=10.200.10.1
  echo "Aliasing host ip: $HOST_IP"
  sudo ifconfig lo0 alias $HOST_IP/24
fi

docker-compose --verbose build --pull --no-cache
docker-compose --verbose up -d
