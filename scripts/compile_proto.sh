#!/bin/bash

set -e

rm -rf toc-trade-protobuf
git clone git@github.com:ToC-Taiwan/toc-trade-protobuf.git

rm -rf pb
mkdir pb

protoc \
    --go_out=. \
    --go-grpc_out=. \
    --proto_path=./toc-trade-protobuf/protos/v3/app \
    --proto_path=./toc-trade-protobuf/protos/v3/forwarder \
    ./toc-trade-protobuf/protos/v3/*/*.proto

rm -rf toc-trade-protobuf
git add ./pb
