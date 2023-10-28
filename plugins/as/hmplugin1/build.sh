#!/bin/bash
pushd `dirname $0` > /dev/null

npm install
npm run asbuild

popd > /dev/null
