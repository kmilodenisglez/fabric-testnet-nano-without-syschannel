#!/usr/bin/env sh
set -eu

docker stop $(docker ps -aqf "NAME=explorer*"); docker rm -f $(docker ps -aqf "NAME=explorer*"); docker volume prune; docker network prune;