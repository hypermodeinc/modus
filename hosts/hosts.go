package hosts

import (
	"context"
	"fmt"
	"os"
	"strings"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/manifestdata"

	"github.com/hypermodeAI/manifest"
)

const hostKeyPrefix = "HYP_HOST_KEY_"

const HypermodeHost string = "hypermode"
const OpenAIHost string = "openai"

func GetHost(hostName string) (manifest.HostInfo, error) {

	if hostName == HypermodeHost {
		return manifest.HostInfo{Name: HypermodeHost}, nil
	}

	for _, host := range manifestdata.Manifest.Hosts {
		if host.Name == hostName {
			return host, nil
		}
	}

	return manifest.HostInfo{}, fmt.Errorf("a host '%s' was not found", hostName)
}

func GetHostForUrl(url string) (manifest.HostInfo, error) {
	for _, host := range manifestdata.Manifest.Hosts {
		// case insensitive version of strings.HasPrefix
		if len(url) >= len(host.Endpoint) && strings.EqualFold(host.Endpoint, url[:len(host.Endpoint)]) {
			return host, nil
		}
	}

	return manifest.HostInfo{}, fmt.Errorf("a host for url '%s' was not found in the manifest", url)
}

func GetHostKey(ctx context.Context, host manifest.HostInfo) (string, error) {
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

		keyEnvVar := hostKeyPrefix + strings.ToUpper(strings.ReplaceAll(host.Name, "-", "_"))
		key = os.Getenv(keyEnvVar)
		if key != "" {
			return key, nil
		} else {
			err = fmt.Errorf("environment variable '%s' not found", keyEnvVar)
		}
	}

	return "", fmt.Errorf("error getting key for host '%s': %w", host.Name, err)
}

func getWellKnownEnvironmentVariable(host manifest.HostInfo) string {

	// Some model hosts have well-known environment variables that are used to store the model key.
	// We should support these to make it easier for users to set up their environment.
	// We can expand this list as we add more model hosts.

	switch host.Name {
	case OpenAIHost:
		return os.Getenv("OPENAI_API_KEY")
	}
	return ""
}
