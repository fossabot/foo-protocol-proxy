#!/usr/bin/env bash

set -eoux

CURRENT_DIR=$(dirname "$0")

. ${CURRENT_DIR}/functions.sh
. ${CURRENT_DIR}/common.sh

while getopts 'f:l:h:r:p:t:' name
do
    case $name in
    f)  FORWARDING_PORT="$OPTARG";;
    l)  LISTENING_PORT="$OPTARG";;
    h)  HEALTH_PORT="$OPTARG";;
    p)  HTTP_PORT="$OPTARG";;
    r)  RECOVERY_PATH="$OPTARG";;
    ?)  printf 'Usage: %s [-f forwarding port] [-l listening port] [-p http port] [-h health port] [-r recovery path]\n' $0
          exit 2;;
    esac
done

cat ${CURRENT_DIR}/../templates/env/app.env.dist > $ENV_FILE_PATH

replaceInFile '{{FORWARDING_PORT}}' $FORWARDING_PORT $ENV_FILE_PATH
replaceInFile '{{LISTENING_PORT}}' $LISTENING_PORT $ENV_FILE_PATH
replaceInFile '{{HEALTH_PORT}}' $HEALTH_PORT $ENV_FILE_PATH
replaceInFile '{{HTTP_PORT}}' $HTTP_PORT $ENV_FILE_PATH
replaceInFile '{{RECOVERY_PATH}}' $RECOVERY_PATH $ENV_FILE_PATH
