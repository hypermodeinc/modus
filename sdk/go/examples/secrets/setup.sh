#!/bin/bash

set -ex

KIND_VERSION=v0.29.0
KIND_CLUSTER=modus
KIND_CONF_FILE="cluster-config.yaml"

KUBECTL_VERSION=v1.32.3
KUBECTL_PATH=/usr/local/bin/kubectl

set +e
KIND=$(command -v kind)
KUBECTL=$(command -v kubectl)
set -e

# install kind
if [[ -z ${KIND} ]]; then
	go install sigs.k8s.io/kind@"${KIND_VERSION}"
	KIND=$(command -v kind)
fi

# install kubectl
if [[ -z ${KUBECTL} ]]; then
	curl -OL "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/$(go env GOARCH)/kubectl"
	chmod +x ./kubectl
	mv ./kubectl "${KUBECTL_PATH}"
	KUBECTL=$(command -v kubectl)
fi

# setup kind cluster configuration
cat >"${KIND_CONF_FILE}" <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
EOF

# create kind cluster
${KIND} create cluster --name "${KIND_CLUSTER}" --config "${KIND_CONF_FILE}"

# delete conf file
rm -f "${KIND_CONF_FILE}"

# get cluster info
${KUBECTL} cluster-info

# create a secret in the default namespace with foo: bar as the value
${KUBECTL} create secret -n default generic example --from-literal=foo=bar
