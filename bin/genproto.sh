#!/usr/bin/env bash

set -e

export ROOT=`git rev-parse --show-toplevel`

protoc -I "${ROOT}" --go_out "${ROOT}/server" "${ROOT}/proto/colony.proto"

