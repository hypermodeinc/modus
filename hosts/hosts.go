package hosts

import (
	"context"
	"fmt"
	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/manifest"
	"os"
	"strings"
)

const hostKeyPrefix = "HYP_HOST_KEY_"
const OpenAIHost string = "openai"

func GetHost(hostName string) (manifest.Host, error) {
	for _, host := range manifest.HypermodeData.Hosts {
		if host.Name == hostName {
			return host, nil
		}
	}

	return manifest.Host{}, fmt.Errorf("a host '%s' was not found", hostName)
}

func GetHostKey(ctx context.Context, host manifest.Host) (string, error) {
	var key string
	var err error

	if config.UseAwsSecrets {
		var hostKey string
		ns := os.Getenv("NAMESPACE")
		if ns == "" {
			hostKey = host.Name
		} else {
			hostKey = ns + "/" + host.Name
		}
		// Get the model key from AWS Secrets Manager, using the model name as the secret.
		key, err = aws.GetSecretString(ctx, hostKey)
		if key != "" {
			return key, nil
		}
	} else {
		// Try well-known environment variables first, then model-specific environment variables.
		key := getWellKnownEnvironmentVariable(host)
		if key != "" {
			return key, nil
		}

		keyEnvVar := hostKeyPrefix + strings.ToUpper(host.Name)
		key = os.Getenv(keyEnvVar)
		if key != "" {
			return key, nil
		} else {
			err = fmt.Errorf("environment variable '%s' not found", keyEnvVar)
		}
	}

	return "", fmt.Errorf("error getting key for model '%s': %w", host.Name, err)
}

func getWellKnownEnvironmentVariable(host manifest.Host) string {

	// Some model hosts have well-known environment variables that are used to store the model key.
	// We should support these to make it easier for users to set up their environment.
	// We can expand this list as we add more model hosts.

	switch host.Name {
	case OpenAIHost:
		return os.Getenv("OPENAI_API_KEY")
	}
	return ""
}
