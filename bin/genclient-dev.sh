#!/usr/bin/env bash

set -e

export ROOT=`git rev-parse --show-toplevel`

cd "${ROOT}/client"
go-bindata -debug -o "${ROOT}/server/client.go" -prefix "resources/public/" resources/public/...
