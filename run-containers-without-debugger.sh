#!/bin/bash

echo "terminating any running containers.."
docker-compose down
docker-compose -f docker-compose-delve.yml down

echo "creating images and starting containers..."
docker-compose up -d --build
