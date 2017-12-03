#!/usr/bin/env bash

set -eux

case "${TRAVIS_OS_NAME}" in
    linux)
        sudo apt-get update -qq && sudo apt-get install -qy python-software-properties
        ;;
esac
