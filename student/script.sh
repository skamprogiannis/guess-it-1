#!/bin/sh
: "${GOCACHE:=/tmp/guess-it-go-build-cache}"
export GOCACHE

cd "$(dirname "$0")" || exit 1
go run .
