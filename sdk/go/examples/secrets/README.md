# Secrets

## Setup

```bash
bash setup.sh
bash build.sh
modus_runtime -appPath ./build -useKubernetesSecret -kubernetesSecretName default/example
```

## Testing

1. Call the GetSecretValue() function with "foo" as the input argument
2. It should return "bar"

## Cleanup

```bash
bash teardown.sh
```
