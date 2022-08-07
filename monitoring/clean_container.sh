#!/usr/bin/env sh
set -eu

docker stop $(docker ps -aqf "NAME=prometheus*"); docker rm -f $(docker ps -aqf "NAME=prometheus*"); docker volume prune; docker network prune;
docker stop $(docker ps -aqf "NAME=grafana*"); docker rm -f $(docker ps -aqf "NAME=grafana*"); docker volume prune; docker network prune;