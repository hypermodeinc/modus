#!/usr/bin/env bash
# Runs "go mod tidy" on all projects.

set -euo pipefail
trap 'cd "${PWD}"' EXIT
cd "$(dirname "$0")"
cd ..

echo "Tidying ."
go mod tidy

cd examples
for example in *; do
	if [[ -d ${example} ]]; then
		cd "${example}"
		echo "Tidying examples/${example}"
		go mod tidy
		cd ..
	fi
done
