#!/bin/bash

rm -rf protos
mkdir -p protos

protoc \
  --go_out=. \
  --go_opt=Mhm.proto=hmruntime/protos \
  --go_opt=module=hmruntime \
  --proto_path=../protos \
  hm.proto

go build
