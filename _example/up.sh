#!/bin/bash

# The up.sh script is used to start the services as an example.
# This example spawns some instances, 5 by default, of the imhost service by
# using the --scale flag.

set -eu

num_instances=${1:-5}

echo "- docker compose up: starting ${num_instances} instances of imhost service."
docker compose up \
    --quiet-pull \
    --force-recreate \
    --remove-orphans \
    --detach \
    --scale imhost="${num_instances}"

# Defer cleaning up the docker compose services.
cleanup() {
    docker compose down 1>/dev/null 2>/dev/null && echo "- docker compose down: done."
}
trap cleanup EXIT

sleep 3

echo "- stress test: starting stress test with ${num_instances} instances."
go run ./stress/main.go -s "${num_instances}"
