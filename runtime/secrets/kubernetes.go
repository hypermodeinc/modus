/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package secrets

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/zerologr"
	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/sentryutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	pkgcache "sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
)

var cacheSyncTimeout = 10 * time.Second

type kubernetesSecretsProvider struct {
	k8sClient       client.Client
	secretName      string
	secretNamespace string
}

func (sp *kubernetesSecretsProvider) initialize(ctx context.Context) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	fullName := app.Config().KubernetesSecretName()
	if fullName == "" {
		const msg = "A Kubernetes secret name is required when using Kubernetes secrets. Exiting."
		sentryutils.CaptureError(ctx, nil, msg)
		logger.Fatal(ctx).Msg(msg)
	}

	parts := strings.SplitN(fullName, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		const msg = "Kubernetes secret name must be in the format <namespace>/<name>"
		sentryutils.CaptureError(ctx, nil, msg, sentryutils.WithData("input", fullName))
		logger.Fatal(ctx).Str("input", fullName).Msg(msg)
	}
	sp.secretNamespace = parts[0]
	sp.secretName = parts[1]

	// Important to hook up the logger so we see any logs from the controller-runtime package.
	ctrl.SetLogger(zerologr.New(logger.Get(ctx)))

	cli, cache, err := newK8sClientForSecret(sp.secretNamespace, sp.secretName)
	if err != nil {
		const msg = "Failed to initialize Kubernetes client."
		sentryutils.CaptureError(ctx, err, msg,
			sentryutils.WithData("namespace", sp.secretNamespace),
			sentryutils.WithData("name", sp.secretName))
		logger.Fatal(ctx, err).Msg(msg)
	}
	sp.k8sClient = cli

	// Start the cache informer in a separate goroutine.
	// This will open a watch connection to the Kubernetes API server,
	// and receive a stream of changes from the API server.
	// See: https://kubernetes.io/docs/reference/using-api/api-concepts/#efficient-detection-of-changes
	//
	// This cache is then used by k8sClient.Get(), so that it does not need
	// to call the API server every time.
	go func() {
		if err := cache.Start(ctx); err != nil {
			const msg = "Failed to start Kubernetes client cache."
			sentryutils.CaptureError(ctx, err, msg,
				sentryutils.WithData("namespace", sp.secretNamespace),
				sentryutils.WithData("name", sp.secretName))
			logger.Fatal(ctx, err).Msg(msg)
		}
		logger.Info(ctx).Msg("Kubernetes client cache stopped.")
	}()

	// Wait for initial cache sync - but not indefinitely.
	waitCtx, cancel := context.WithTimeout(ctx, cacheSyncTimeout)
	defer cancel()
	if !cache.WaitForCacheSync(waitCtx) {
		const msg = "Failed to sync Kubernetes client cache."
		sentryutils.CaptureError(ctx, nil, msg,
			sentryutils.WithData("namespace", sp.secretNamespace),
			sentryutils.WithData("name", sp.secretName))
		logger.Fatal(ctx).Msg(msg)
	}

	logger.Info(ctx).
		Str("secret", sp.secretName).
		Str("namespace", sp.secretNamespace).
		Msg("Kubernetes secrets provider initialized.")
}

func (sp *kubernetesSecretsProvider) getConnectionSecrets(ctx context.Context, connection manifest.ConnectionInfo) (map[string]string, error) {
	secrets := make(map[string]string)

	if sp.k8sClient == nil || sp.secretName == "" {
		return nil, fmt.Errorf("Kubernetes client or secret name not initialized")
	}

	secret := &corev1.Secret{}
	err := sp.k8sClient.Get(ctx, client.ObjectKey{
		Name:      sp.secretName,
		Namespace: sp.secretNamespace,
	}, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes secret: %w", err)
	}

	prefix := "MODUS_" + strings.ToUpper(strings.ReplaceAll(connection.ConnectionName(), "-", "_")) + "_"
	for k, v := range secret.Data {
		if strings.HasPrefix(k, prefix) {
			secrets[k[len(prefix):]] = string(v)
		}
	}
	return secrets, nil
}

func (sp *kubernetesSecretsProvider) getSecretValue(ctx context.Context, name string) (string, error) {
	if sp.k8sClient == nil || sp.secretName == "" {
		return "", fmt.Errorf("Kubernetes client or secret name not initialized")
	}
	secret := &corev1.Secret{}
	err := sp.k8sClient.Get(ctx, client.ObjectKey{
		Name:      sp.secretName,
		Namespace: sp.secretNamespace,
	}, secret)
	if err != nil {
		return "", fmt.Errorf("failed to get kubernetes secret: %w", err)
	}
	v, ok := secret.Data[name]
	if !ok {
		return "", fmt.Errorf("secret %s was not found", name)
	}
	return string(v), nil
}

func (sp *kubernetesSecretsProvider) hasSecret(ctx context.Context, name string) bool {
	if sp.k8sClient == nil || sp.secretName == "" {
		return false
	}
	secret := &corev1.Secret{}
	err := sp.k8sClient.Get(ctx, client.ObjectKey{
		Name:      sp.secretName,
		Namespace: sp.secretNamespace,
	}, secret)
	if err != nil {
		return false
	}
	_, ok := secret.Data[name]
	return ok
}

// newK8sClientForSecret creates a new Kubernetes client that watches
// only a single secret `name` in namespace `ns`.
func newK8sClientForSecret(ns, name string) (client.Client, pkgcache.Cache, error) {
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	cfg, err := ctrl.GetConfig()
	if err != nil {
		return nil, nil, err
	}

	cl, err := cluster.New(cfg, func(clusterOptions *cluster.Options) {
		clusterOptions.Scheme = scheme

		// Only cache the one secret we need, so we don't waste memory space
		clusterOptions.Cache.DefaultNamespaces = map[string]pkgcache.Config{
			ns: {
				FieldSelector: fields.OneTermEqualSelector(metav1.ObjectNameField, name),
			},
		}
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to construct cluster: %v", err)
	}

	return cl.GetClient(), cl.GetCache(), nil
}
