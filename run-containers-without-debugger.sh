#!/bin/bash

echo "terminating any running containers.."
docker-compose down

echo "creating images and starting containers..."
docker-compose up -d --build
