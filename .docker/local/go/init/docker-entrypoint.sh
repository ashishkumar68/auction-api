#!/bin/bash

set -e

sleep 5

go mod vendor

exec "$@"