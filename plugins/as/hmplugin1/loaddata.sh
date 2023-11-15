#!/bin/bash
set -e
pushd `dirname $0` > /dev/null

echo "Writing database schema..."
curl -X POST http://localhost:8080/admin/schema \
  -H 'Content-Type: application/graphql' \
  --data-binary '@schema.graphql' \
  -s -o /dev/null

echo "Loading sample data..."
curl -X POST http://localhost:8080/graphql \
  -H 'Content-Type: application/graphql' \
  --data-binary '@sampledata.graphql' \
  -s -o /dev/null

echo "Done!"
popd > /dev/null
