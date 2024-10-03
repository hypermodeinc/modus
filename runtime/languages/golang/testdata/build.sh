#!/bin/bash
PROJECTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd ../../../../sdk/go/tools/hypbuild > /dev/null
go run . "$PROJECTDIR"
popd > /dev/null
