#!/bin/bash

set -e

sleep 3

go mod vendor

exec "$@"