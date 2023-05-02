#!/bin/bash

proto_gen_path="pkg"
proto_path="./proto"

mkdir -p ./$proto_gen_path

protoc --go_out=$proto_gen_path --go_opt=paths=source_relative \
  --go-grpc_out=$proto_gen_path --go-grpc_opt=paths=source_relative $proto_path/*

if [ $? -eq 0 ]
  then echo "proto generated in $proto_gen_path"
fi
