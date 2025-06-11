# Secrets

## Setup

This example requires a Kubernetes cluster to be running. To simplify this, we use
[kind](https://kind.sigs.k8s.io/) to create a local cluster.

```bash
# Setup local kind cluster for testing
bash setup.sh

# Build the Modus app
bash build.sh

# Run the Modus app
modus_runtime -appPath ./build -useKubernetesSecret -kubernetesSecretName default/example
```

## Testing

1. Call the GetSecretValue() function with "foo" as the input argument
2. It should return "bar"

## Cleanup

```bash
# Delete the kind cluster
bash teardown.sh
```
