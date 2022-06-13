#!/bin/bash

git clone git@gitlab.tocraw.com:root/toc-trade-protobuf.git
protoc --proto_path=./toc-trade-protobuf --go_out=. --go-grpc_out=. ./toc-trade-protobuf/src/*.proto
rm -rf toc-trade-protobuf

git add ./pb/sinopac_forwarder_grpc.pb.go
git add ./pb/sinopac_forwarder.pb.go
