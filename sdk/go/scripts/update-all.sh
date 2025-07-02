#!/usr/bin/env bash
# Runs "go mod tidy" on all projects.

set -euo pipefail
trap 'cd "${PWD}"' EXIT
cd "$(dirname "$0")"
cd ..

echo "Tidying ."
go mod tidy

cd templates
for template in *; do
	if [[ -d ${template} ]]; then
		cd "${template}"
		echo "Updating templates/${template}"
		go get -u
		go mod tidy
		cd ..
	fi
done

cd ../examples
for example in *; do
	if [[ -d ${example} ]]; then
		cd "${example}"
		echo "Updating examples/${example}"
		go get -u
		go mod tidy
		cd ..
	fi
done
