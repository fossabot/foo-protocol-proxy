#!/usr/bin/env bash

set -eoux

CURRENT_DIR=$(dirname "$0")

# .env file replacements:
. ${CURRENT_DIR}/prepare-env-file.sh

# 🐳 Dockerfile replacements:
. ${CURRENT_DIR}/prepare-docker-file.sh

# 🐳 Docker compose file replacements:
. ${CURRENT_DIR}/prepare-docker-compose-file.sh
