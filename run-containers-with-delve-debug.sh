#!/bin/bash

echo "terminating any running containers.."
docker-compose down

echo "building images..."
docker-compose build --build-arg AIR_FILE_NAME=".air.delve-debug.toml"

echo "created images, starting containers..."
docker-compose up -d
