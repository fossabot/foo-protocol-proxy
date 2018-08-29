#!/usr/bin/env bash

set -eoux

CURRENT_DIR=$(dirname "$0")
ENV_FILE='.env'

# Environment variables.
FORWARDING_PORT=8001
LISTENING_PORT=8002
HTTP_PORT=8088
RECOVERY_PATH='.data/recovery.json'

while getopts 'f:l:h:r:p:t:' name
do
    case $name in
    f)  FORWARDING_PORT="$OPTARG";;
    l)  LISTENING_PORT="$OPTARG";;
    h)  HTTP_PORT="$OPTARG";;
    r)  RECOVERY_PATH="$OPTARG";;
    ?)  printf "Usage: %f: [-f forwarding port] [-l listening port] [-h http address] [-r recovery path]\n" $0
          exit 2;;
    esac
done

. ${CURRENT_DIR}/includes.sh

cat ${CURRENT_DIR}/../templates/env/app.env.dist > $ENV_FILE

replaceInFile '{{FORWARDING_PORT}}' $FORWARDING_PORT $ENV_FILE
replaceInFile '{{LISTENING_PORT}}' $LISTENING_PORT $ENV_FILE
replaceInFile '{{HTTP_PORT}}' $HTTP_PORT $ENV_FILE
replaceInFile '{{RECOVERY_PATH}}' $RECOVERY_PATH $ENV_FILE
