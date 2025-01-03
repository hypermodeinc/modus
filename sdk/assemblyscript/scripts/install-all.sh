#!/usr/bin/env bash
# Runs "npm install" on all projects.

set -euo pipefail
trap "cd \"${PWD}\"" EXIT
cd "$(dirname "$0")"
cd ..

cd src
npm install

cd ../examples
for example in *; do
	if [[ -d ${example} ]]; then
		cd "${example}"
		npm install
		cd ..
	fi
done
