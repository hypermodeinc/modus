#!/bin/bash
pushd `dirname $0` > /dev/null

npm install

rm -rf protos
mkdir -p protos

protoc \
  --plugin=protoc-gen-as=node_modules/.bin/as-proto-gen \
  --as_out=protos \
  --proto_path=../../../protos \
  hm.proto

npm run asbuild

popd > /dev/null
