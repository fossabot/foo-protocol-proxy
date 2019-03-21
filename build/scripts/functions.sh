#!/usr/bin/env bash

set -eoux

# Gets the current platform.
getPlatform() {
  echo $(uname -s | tr '[:upper:]' '[:lower:]') 2>&1
}

# Checks if the host is Linux.
isDarwinHost() {
  local platform=$1

  if [[ $platform == 'darwin' ]]; then
     return 0
  fi

  return 1
}

replaceInFile() {
    local search=$1
    local replace=$2
    local filePath=$3

    if ( isDarwinHost $(getPlatform) )
    then
        # http://lists.gnu.org/archive/html/bug-gnu-utils/2013-09/msg00000.html
        sed -i '' -e "s|${search}|${replace}|g" $filePath 2>&1
    else
        sed -i -e "s|${search}|${replace}|g" $filePath 2>&1
    fi
}

# Should an ip alias used based on the platform.
if ( isDarwinHost $(getPlatform) )
then
  export HOST_IP=10.200.10.1
  echo "Aliasing host ip: $HOST_IP"
  sudo ifconfig lo0 alias $HOST_IP/24
fi
