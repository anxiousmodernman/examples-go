#!/bin/bash

set -e

SRC_DIR="proto"
DEST_DIR="server"

protoc_paths="-I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis"

protoc ${protoc_paths} --go_out=plugins=grpc:. ${SRC_DIR}/server.proto

go build
