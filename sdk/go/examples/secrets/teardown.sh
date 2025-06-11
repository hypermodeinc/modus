#!/bin/bash

set -ex

KIND_CLUSTER=modus
KIND=$(command -v kind)

# delete kind cluster
${KIND} delete cluster --name "${KIND_CLUSTER}"
