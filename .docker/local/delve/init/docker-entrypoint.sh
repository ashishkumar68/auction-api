#!/bin/bash

set -e

go build -gcflags="all=-N -l" -o ${GOBIN}/server

exec "$@"