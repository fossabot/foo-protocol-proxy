#!/usr/bin/env sh

TAG="0.0.1"
FORWARDING_PORT=8001
LISTENING_PORT=8002
HTTP_PORT=8088
RECOVERY_PATH="data/recovery.json"

while getopts "f:l:h:r:" name
do
    case $name in
    f)  FORWARDING_PORT="$OPTARG";;
    l)  LISTENING_PORT="$OPTARG";;
    h)  HTTP_PORT="$OPTARG";;
    r)  RECOVERY_PATH="$OPTARG";;
    ?)  printf "Usage: %f: [-f forwarding port]\n" $0
          exit 2;;
    esac
done

export TAG="$TAG"
export FORWARDING_PORT="$FORWARDING_PORT"
export LISTENING_PORT="$LISTENING_PORT"
export HTTP_PORT="$HTTP_PORT"
export RECOVERY_PATH="$RECOVERY_PATH"

# Gets the current platform.
getPlatform() {
  echo $(uname -s)
}

# Should an ip alias used based on the platform.
shouldUseIpAlias() {
  local platform=$1
  if [[ $platform == 'linux' ]]; then
     return 1
  elif [[ $platform == 'Darwin' ]]; then
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
