package envfiles

import (
	"context"
	"sync"

	"github.com/hypermodeinc/modus/runtime/app"
)

type EnvFilesLoadedCallback = func(ctx context.Context)

var manifestLoadedCallbacks []EnvFilesLoadedCallback
var eventsMutex = sync.RWMutex{}

func RegisterEnvFilesLoadedCallback(callback EnvFilesLoadedCallback) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()
	manifestLoadedCallbacks = append(manifestLoadedCallbacks, callback)
}

func triggerEnvFilesLoaded(ctx context.Context) error {
	if app.IsShuttingDown() {
		return nil
	}

	eventsMutex.RLock()
	defer eventsMutex.RUnlock()

	for _, callback := range manifestLoadedCallbacks {
		callback(ctx)
	}
	return nil
}
