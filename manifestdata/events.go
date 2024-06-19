/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifestdata

import (
	"context"

	"github.com/hypermodeAI/manifest"
)

type ManifestLoadedCallback = func(ctx context.Context, manifest manifest.HypermodeManifest) error

var manifestLoadedCallbacks []ManifestLoadedCallback

func RegisterManifestLoadedCallback(callback ManifestLoadedCallback) {
	manifestLoadedCallbacks = append(manifestLoadedCallbacks, callback)
}

func triggerManifestLoaded(ctx context.Context, manifest manifest.HypermodeManifest) error {
	for _, callback := range manifestLoadedCallbacks {
		err := callback(ctx, manifest)
		if err != nil {
			return err
		}
	}
	return nil
}
