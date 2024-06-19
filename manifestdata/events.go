/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifestdata

import (
	"context"
)

type ManifestLoadedCallback = func(ctx context.Context) error

var manifestLoadedCallbacks []ManifestLoadedCallback

func RegisterManifestLoadedCallback(callback ManifestLoadedCallback) {
	manifestLoadedCallbacks = append(manifestLoadedCallbacks, callback)
}

func triggerManifestLoaded(ctx context.Context) error {
	for _, callback := range manifestLoadedCallbacks {
		err := callback(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
