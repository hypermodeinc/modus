package functions

import "context"

type FunctionsLoadedCallback = func(ctx context.Context) error

var functionsLoadedCallbacks []FunctionsLoadedCallback

func RegisterFunctionsLoadedCallback(callback FunctionsLoadedCallback) {
	functionsLoadedCallbacks = append(functionsLoadedCallbacks, callback)
}

func triggerFunctionsLoaded(ctx context.Context) error {
	for _, callback := range functionsLoadedCallbacks {
		err := callback(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
