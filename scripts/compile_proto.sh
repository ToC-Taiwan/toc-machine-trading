#!/bin/bash

git clone git@gitlab.tocraw.com:root/toc-trade-protobuf.git
rm -rf pb
protoc --proto_path=./toc-trade-protobuf --go_out=. --go-grpc_out=. ./toc-trade-protobuf/pb/*.proto
rm -rf toc-trade-protobuf
git add ./pb
