/*
 * Copyright 2024 Hypermode, Inc.
 */

package secrets

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/hypermodeAI/manifest"
)

var provider secretsProvider

type secretsProvider interface {
	initialize(ctx context.Context)
	getSecretValue(name string) (string, error)
	getHostSecrets(host manifest.HostInfo) (map[string]string, error)
}

func Initialize(ctx context.Context) {
	if config.UseAwsSecrets {
		provider = &awsSecretsProvider{}
	} else {
		provider = &localSecretsProvider{}
	}

	provider.initialize(ctx)
}

func GetSecretValue(name string) (string, error) {
	return provider.getSecretValue(name)
}

func GetHostSecrets(host manifest.HostInfo) (map[string]string, error) {
	return provider.getHostSecrets(host)
}

func GetHostSecret(host manifest.HostInfo, secretName string) (string, error) {
	secrets, err := GetHostSecrets(host)
	if err != nil {
		return "", err
	}

	if val, ok := secrets[secretName]; ok {
		return val, nil
	}

	return "", fmt.Errorf("could not find secret '%s' for host '%s'", secretName, host.HostName())
}

// ApplyHostSecretsToHttpRequest evaluates the given request and replaces any placeholders
// present in the query parameters and headers with their secret values for the given host.
func ApplyHostSecretsToHttpRequest(ctx context.Context, host *manifest.HTTPHostInfo, req *http.Request) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	// get secrets for the host
	secrets, err := GetHostSecrets(host)
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

// ApplyHostSecretsToString evaluates the given string and replaces any placeholders
// present in the string with their secret values for the given host.
func ApplyHostSecretsToString(ctx context.Context, host manifest.HostInfo, str string) (string, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	secrets, err := GetHostSecrets(host)
	if err != nil {
		return "", err
	}
	return applySecretsToString(ctx, secrets, str), nil
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
			if ok {
				return val
			}

			logger.Warn(ctx).Str("name", nameKey).Msg("Secret not found.")
			return ""
		}

		// base64 secret template
		userKey := submatches[1]
		passKey := submatches[2]
		if userKey != "" && passKey != "" {
			user, userOk := secrets[userKey]
			if !userOk {
				logger.Warn(ctx).Str("name", userKey).Msg("Secret not found.")
			}

			pass, passOk := secrets[passKey]
			if !passOk {
				logger.Warn(ctx).Str("name", passKey).Msg("Secret not found.")
			}

			if userOk && passOk {
				return base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
			} else {
				return ""
			}
		}

		logger.Warn(ctx).Str("template", match).Msg("Invalid secret variable template.")
		return match
	})
}
