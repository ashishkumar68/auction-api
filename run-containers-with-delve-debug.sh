#!/bin/bash

echo "terminating any running containers.."
docker-compose down
docker-compose -f docker-compose-delve.yml down

echo "building images..."
docker-compose -f docker-compose-delve.yml build

echo "created images, starting containers..."
docker-compose -f docker-compose-delve.yml up -d
