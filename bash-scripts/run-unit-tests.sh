#!/bin/bash

echo "Running unit tests.."
go test -race ./... -p 1