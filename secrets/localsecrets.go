/*
 * Copyright 2024 Hypermode, Inc.
 */

package secrets

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hypermodeAI/manifest"
)

type localSecretsProvider struct {
}

func (sp *localSecretsProvider) initialize(ctx context.Context) {
}

func (sp *localSecretsProvider) getHostSecrets(ctx context.Context, host manifest.HostInfo) (map[string]string, error) {
	prefix := "HYPERMODE_" + strings.ToUpper(strings.ReplaceAll(host.HostName(), "-", "_")) + "_"
	secrets := make(map[string]string)
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, prefix) {
			pair := strings.SplitN(e, "=", 2)
			secrets[pair[0][len(prefix):]] = pair[1]
		}
	}

	return secrets, nil
}

func (sp *localSecretsProvider) getSecretValue(ctx context.Context, name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("environment variable %s was not found", name)
	}

	return v, nil
}
