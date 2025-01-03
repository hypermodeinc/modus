#!/bin/bash

# This build script works best for examples that are in this repository.
# If you are using this as a template for your own project, you may need to modify this script,
# to invoke the modus-go-build tool with the correct path to your project.

PROJECTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
pushd ../../tools/modus-go-build >/dev/null || exit
go run . "${PROJECTDIR}"
exit_code=$?
popd >/dev/null || exit
exit "${exit_code}"
