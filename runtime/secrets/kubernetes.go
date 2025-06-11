/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package secrets

import (
	"context"
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

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

type kubernetesSecretsProvider struct {
	k8sClient       client.Client
	secretName      string
	secretNamespace string
}

func (sp *kubernetesSecretsProvider) initialize(ctx context.Context) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	fullName := app.Config().KubernetesSecretName()
	if fullName == "" {
		logger.Fatal(ctx).Msg("A Kubernetes secret name is required when using Kubernetes secrets. Exiting.")
	}

	parts := strings.SplitN(fullName, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		logger.Fatal(ctx).Str("input", fullName).Msg("Kubernetes secret name must be in the format <namespace>/<name>")
	}
	sp.secretNamespace = parts[0]
	sp.secretName = parts[1]

	cli, cache, err := newK8sClientForSecret(sp.secretNamespace, sp.secretName)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to initialize Kubernetes client.")
	}
	sp.k8sClient = cli

	// Start the cache informer in a separate goroutine.
	// This will open a watch connection to the Kubernetes API server,
	// and receive a stream of changes from the API server.
	// See: https://kubernetes.io/docs/reference/using-api/api-concepts/#efficient-detection-of-changes
	//
	// This cache is then used by k8sClient.Get(), so that it does not need
	// to call the API server everytime.
	go func() {
		if err := cache.Start(ctx); err != nil {
			logger.Fatal(ctx).Err(err).Msg("Failed to start Kubernetes client cache.")
		}
		logger.Info(ctx).Msg("Kubernetes client cache stopped.")
	}()

	// Wait for initial cache sync
	if !cache.WaitForCacheSync(ctx) {
		logger.Fatal(ctx).Msg("Failed to sync Kubernetes client cache.")
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

	cl, err := cluster.New(ctrl.GetConfigOrDie(), func(clusterOptions *cluster.Options) {
		clusterOptions.Scheme = scheme
		clusterOptions.Cache.DefaultNamespaces = make(map[string]pkgcache.Config)
		// Only cache the one secret we need, so we don't waste memory space
		clusterOptions.Cache.DefaultNamespaces[ns] = pkgcache.Config{
			FieldSelector: fields.OneTermEqualSelector(metav1.ObjectNameField, name),
		}
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to construct cluster: %v", err)
	}

	return cl.GetClient(), cl.GetCache(), nil
}
