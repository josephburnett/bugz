#!/usr/bin/env bash

set -e

export ROOT=`git rev-parse --show-toplevel`

cd "${ROOT}"
# bin/genproto.sh
bin/genclient.sh

cd "${ROOT}/server"
go build -o "${ROOT}/build/colony" *.go
