#!/usr/bin/env bash
# Runs "npm install" on all projects.

set -euo pipefail
trap 'cd "${PWD}"' EXIT
cd "$(dirname "$0")"
cd ..

cd src
npm install
npm run build:transform

cd ../templates
for template in *; do
	if [[ -d ${template} ]]; then
		cd "${template}"
		npm install
		npm run build
		cd ..
	fi
done

cd ../examples
for example in *; do
	if [[ -d ${example} ]]; then
		cd "${example}"
		npm install
		npm run build
		cd ..
	fi
done
