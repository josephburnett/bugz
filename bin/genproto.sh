#!/usr/bin/env bash

set -e

export ROOT=`git rev-parse --show-toplevel`

echo "Generating protobuf Go client..."
protoc -I "${ROOT}" --go_out "${ROOT}/server" "${ROOT}/proto/view/view.proto"
protoc -I "${ROOT}" --go_out "${ROOT}/server" "${ROOT}/proto/event/event.proto"

echo "Generating protobuf JavaScript client..."
cd "${ROOT}"
protoc -I "${ROOT}" --"js_out=library=client/proto/proto_libs,binary:." \
       "${ROOT}/proto/view/view.proto" \
       "${ROOT}/proto/event/event.proto"

