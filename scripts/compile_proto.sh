#!/bin/bash

rm -rf toc-trade-protobuf
git clone git@github.com:ToC-Taiwan/toc-trade-protobuf.git

rm -rf pb
mkdir pb

protoc --proto_path=./toc-trade-protobuf --go_out=. --go-grpc_out=. ./toc-trade-protobuf/*.proto
rm -rf toc-trade-protobuf
git add ./pb
