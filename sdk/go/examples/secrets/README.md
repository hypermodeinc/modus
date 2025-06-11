# Secrets

## Setup

```bash
bash setup.sh
modus_runtime -appPath ./build -useKubernetesSecret -kubernetesSecretName default/example
```

## Testing

1. Call the GetSecretValue() function with "foo" as the input
2. It should return "bar"

## Cleanup

```bash
bash teardown.sh
```
