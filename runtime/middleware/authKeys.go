/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package middleware

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/hypermodeinc/modus/runtime/logger"
)

var (
	globalAuthKeys *AuthKeys
)

type AuthKeys struct {
	pemPublicKeys  map[string]any
	jwksPublicKeys map[string]any
	mu             sync.RWMutex
	quit           chan struct{}
	done           chan struct{}
}

func newAuthKeys() *AuthKeys {
	return &AuthKeys{
		pemPublicKeys:  make(map[string]any),
		jwksPublicKeys: make(map[string]any),
		quit:           make(chan struct{}),
		done:           make(chan struct{}),
	}
}

func (ak *AuthKeys) setPemPublicKeys(keys map[string]any) {
	ak.mu.Lock()
	defer ak.mu.Unlock()
	ak.pemPublicKeys = keys
}

func (ak *AuthKeys) setJwksPublicKeys(keys map[string]any) {
	ak.mu.Lock()
	defer ak.mu.Unlock()
	ak.jwksPublicKeys = keys
}

func (ak *AuthKeys) getPemPublicKeys() map[string]any {
	ak.mu.RLock()
	defer ak.mu.RUnlock()
	return ak.pemPublicKeys
}

func (ak *AuthKeys) getJwksPublicKeys() map[string]any {
	ak.mu.RLock()
	defer ak.mu.RUnlock()
	return ak.jwksPublicKeys
}

func getJwksRefreshMinutes(ctx context.Context) int {
	refreshTimeStr := os.Getenv("MODUS_JWKS_REFRESH_MINUTES")
	if refreshTimeStr == "" {
		return 1440
	}
	refreshTime, err := strconv.Atoi(refreshTimeStr)
	if err != nil {
		logger.Warn(ctx, err).Msg("Invalid MODUS_JWKS_REFRESH_MINUTES value. Using default value of 1440 minutes.")
		return 1440
	}
	return refreshTime
}

func (ak *AuthKeys) worker(ctx context.Context) {
	defer close(ak.done)
	timer := time.NewTimer(time.Duration(getJwksRefreshMinutes(ctx)) * time.Minute)

	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// refresh JWKS keys
			keysStr := os.Getenv("MODUS_JWKS_ENDPOINTS")
			if keysStr != "" {
				keys, err := jwksEndpointsJsonToKeys(ctx, keysStr)
				if err != nil {
					logger.Warn(ctx, err).Msg("Auth JWKS public keys deserializing error")
				} else {
					ak.setJwksPublicKeys(keys)
				}
			}
			timer.Reset(time.Duration(getJwksRefreshMinutes(ctx)) * time.Minute)
		case <-ak.quit:
			return
		}
	}

}
