#!/bin/bash
PROJECTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
pushd ../../../../sdk/go/tools/modus-go-build >/dev/null || exit
go run . "${PROJECTDIR}"
exit_code=$?
popd >/dev/null || exit
exit "${exit_code}"
