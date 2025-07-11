/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package secrets

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hypermodeinc/modus/lib/manifest"
)

type localSecretsProvider struct {
}

func (sp *localSecretsProvider) initialize(ctx context.Context) {
}

func (sp *localSecretsProvider) getConnectionSecrets(ctx context.Context, connection manifest.ConnectionInfo) (map[string]string, error) {
	prefix := "MODUS_" + strings.ToUpper(strings.ReplaceAll(connection.ConnectionName(), "-", "_")) + "_"
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

func (sp *localSecretsProvider) hasSecret(ctx context.Context, name string) bool {
	return os.Getenv(name) != ""
}
