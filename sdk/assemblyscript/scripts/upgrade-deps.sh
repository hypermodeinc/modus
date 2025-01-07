#!/usr/bin/env bash
# Upgrades all dependencies to the latest minor version, on all projects.
# Requires npm-check-updates (https://www.npmjs.com/package/npm-check-updates).

set -euo pipefail
trap 'cd "${PWD}"' EXIT
cd "$(dirname "$0")"
cd ..

cd src
ncu -u -t minor
npm install

cd ../examples
for example in *; do
	if [[ -d ${example} ]]; then
		cd "${example}"
		ncu -u -t minor
		npm install
		cd ..
	fi
done
