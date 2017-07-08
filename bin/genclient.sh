#!/usr/bin/env bash

set -e

export ROOT=`git rev-parse --show-toplevel`

cd "${ROOT}/client"
lein clean
lein cljsbuild once min
go-bindata -o "${ROOT}/server/client.go" -prefix "resources/public/" resources/public/...

