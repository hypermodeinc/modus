#!/usr/bin/env bash
# Prepares templates for release.

set -euo pipefail
trap "cd \"${PWD}\"" EXIT
cd "$(dirname "$0")"

# get the version number from the command line argument
if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <version>"
  exit 1
fi
version=$1

version=${version#"v"}
echo "Preparing templates for release with version ${version}"

cd ../templates
for template in *; do
  if [ -d "${template}" ]; then
    cd "${template}"

    # update the go.mod file to use the new version
    go get -u github.com/hypermodeinc/modus/sdk/go@v${version}

    cd ..
  fi
done

cd ..
tar -czvf ../../templates_go_v${version}.tar.gz templates
