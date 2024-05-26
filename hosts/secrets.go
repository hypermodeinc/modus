/*
 * Copyright 2024 Hypermode, Inc.
 */

package hosts

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/logger"
	"hmruntime/manifestdata"

	"github.com/hypermodeAI/manifest"
)

func GetHostSecrets(ctx context.Context, host manifest.HostInfo) (map[string]string, error) {
	if config.UseAwsSecrets {

		ns := os.Getenv("NAMESPACE")
		if ns == "" {
			if config.GetEnvironmentName() == "dev" {
				user, err := user.Current()
				if err != nil {
					return nil, fmt.Errorf("could not get current user from the os: %w", err)
				}
				ns = "dev/" + user.Username
			} else {
				return nil, fmt.Errorf("NAMESPACE environment variable is not set")
			}
		}

		prefix := strings.Trim(strings.Join([]string{ns, host.Name}, "/"), "/")
		secrets, err := aws.GetSecrets(ctx, prefix)
		if err != nil {
			return nil, err
		}

		// Migrate old auth header secret to the new location
		// TODO: Remove this when we no longer need to support the old manifest format
		oldAuthHeaderSecret, ok := secrets[""]
		if ok {
			if manifestdata.Manifest.Version == 1 {
				secrets[manifest.V1AuthHeaderVariableName] = oldAuthHeaderSecret
				delete(secrets, "")
				logger.Warn(ctx).Msg("Used deprecated auth header secret.  Please update the manifest to use a template such as {{SECRET_NAME}} and migrate the old secret in Secrets Manager.")
			} else {
				logger.Warn(ctx).Msg("The manifest is current, but the deprecated auth header secret was found.  Please remove the old secret in Secrets Manager.")
			}
		}

		return secrets, nil
	} else {
		prefix := "HYPERMODE_" + strings.ToUpper(strings.ReplaceAll(host.Name, "-", "_")) + "_"
		secrets := make(map[string]string)
		for _, e := range os.Environ() {
			if strings.HasPrefix(e, prefix) {
				pair := strings.SplitN(e, "=", 2)
				secrets[pair[0][len(prefix):]] = pair[1]
			}
		}

		return secrets, nil
	}
}

func GetHostSecret(ctx context.Context, host manifest.HostInfo, secretName string) (string, error) {
	secrets, err := GetHostSecrets(ctx, host)
	if err != nil {
		return "", err
	}

	if val, ok := secrets[secretName]; ok {
		return val, nil
	}

	return "", fmt.Errorf("could not find secret '%s' for host '%s'", secretName, host.Name)
}

func ApplyHostSecrets(ctx context.Context, host manifest.HostInfo, req *http.Request) error {

	// get secrets for the host
	secrets, err := GetHostSecrets(ctx, host)
	if err != nil {
		return err
	}

	// apply query parameters from manifest
	q := req.URL.Query()
	for k, v := range host.QueryParameters {
		q.Add(k, applySecretsToString(ctx, secrets, v))
	}
	req.URL.RawQuery = q.Encode()

	// apply headers from manifest
	for k, v := range host.Headers {
		req.Header.Add(k, applySecretsToString(ctx, secrets, v))
	}

	return nil
}

var templateRegex = regexp.MustCompile(`{{\s*(?:base64\((.+?):(.+?)\)|(.+?))\s*}}`)

func applySecretsToString(ctx context.Context, secrets map[string]string, s string) string {
	return templateRegex.ReplaceAllStringFunc(s, func(match string) string {
		submatches := templateRegex.FindStringSubmatch(match)
		if len(submatches) != 4 {
			logger.Warn(ctx).Str("match", match).Msg("Invalid template.")
			return match
		}

		// standard secret template
		nameKey := submatches[3]
		if nameKey != "" {
			val, ok := secrets[nameKey]
			if !ok {
				logger.Warn(ctx).Str("name", nameKey).Msg("Secret not found.")
				return match
			}
			return val
		}

		// base64 secret template
		userKey := submatches[1]
		passKey := submatches[2]
		if userKey != "" && passKey != "" {
			user, ok := secrets[userKey]
			if !ok {
				logger.Warn(ctx).Str("name", userKey).Msg("Secret not found.")
				return match
			}

			pass, ok := secrets[passKey]
			if !ok {
				logger.Warn(ctx).Str("name", passKey).Msg("Secret not found.")
				return match
			}

			return base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
		}

		logger.Warn(ctx).Str("template", match).Msg("Invalid secret variable template.")
		return match
	})
}
