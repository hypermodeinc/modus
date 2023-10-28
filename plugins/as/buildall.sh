#!/bin/bash
pushd `dirname $0` > /dev/null

./hypermode-as/build.sh
./hmplugin1/build.sh

popd > /dev/null
