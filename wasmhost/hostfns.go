/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"hmruntime/languages"
	"hmruntime/logger"
	"hmruntime/plugins/metadata"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

type hfMessageType int

const (
	hfMessageDetail hfMessageType = iota
	hfMessageStarting
	hfMessageCompleted
	hfMessageCancelled
	hfMessageError
)

type HostFunction struct {
	module          string
	name            string
	function        wasm.GoFunction
	wasmParamTypes  []wasm.ValueType
	wasmResultTypes []wasm.ValueType
	messages        map[hfMessageType]string
	msgDetailFunc   any
}

func (hf *HostFunction) Name() string {
	return hf.module + "." + hf.name
}

func WithStartingMessage(text string) func(*HostFunction) {
	return func(hf *HostFunction) {
		hf.messages[hfMessageStarting] = text
	}
}

func WithCompletedMessage(text string) func(*HostFunction) {
	return func(hf *HostFunction) {
		hf.messages[hfMessageCompleted] = text
	}
}

func WithCancelledMessage(text string) func(*HostFunction) {
	return func(hf *HostFunction) {
		hf.messages[hfMessageCancelled] = text
	}
}

func WithErrorMessage(text string) func(*HostFunction) {
	return func(hf *HostFunction) {
		hf.messages[hfMessageError] = text
	}
}

func WithMessageDetail(fn any) func(*HostFunction) {
	return func(hf *HostFunction) {
		hf.msgDetailFunc = fn
	}
}

func (host *WasmHost) RegisterHostFunction(modName, funcName string, fn any, opts ...func(*HostFunction)) error {
	hf, err := prepareHostFunction(modName, funcName, fn, opts...)
	if err != nil {
		return fmt.Errorf("failed to register host function %s.%s: %w", modName, funcName, err)
	}

	host.hostFunctions = append(host.hostFunctions, hf)
	return nil
}

func prepareHostFunction(modName, funcName string, fn any, opts ...func(*HostFunction)) (*HostFunction, error) {
	rvFunc := reflect.ValueOf(fn)
	if rvFunc.Kind() != reflect.Func {
		return nil, fmt.Errorf("host function %s.%s is not a function type", modName, funcName)
	}

	fnType := rvFunc.Type()
	numParams := fnType.NumIn()
	numResults := fnType.NumOut()

	// Optionally, the first parameter can be a context.
	var hasContextParam bool
	if numParams > 0 && fnType.In(0).Implements(contextType) {
		hasContextParam = true
	}

	// Optionally, the last return value can be an error.
	var hasErrorResult bool
	if numResults > 0 && fnType.Out(numResults-1).Implements(errorType) {
		hasErrorResult = true
	}

	// Multi-value returns are not allowed because we would need a way to encode them for
	// languages that don't naturally support them.
	// TODO: In the future, this could be done by having the SDK generate a struct type for the return value.
	if (hasErrorResult && numResults > 2) || (!hasErrorResult && numResults > 1) {
		return nil, fmt.Errorf("host function %s.%s cannot return multiple data values", modName, funcName)
	}

	// TODO: the following assumes a lot.  We should use the language's type system to determine the encoding
	// but we'll need to do this registration when the plugin is loaded so we know what language is being used.
	// For now, we'll assume single-value, one-to-one mapping between Go and Wasm types.

	paramTypes := make([]wasm.ValueType, 0, numParams)
	for i := 0; i < numParams; i++ {
		if hasContextParam && i == 0 {
			continue
		}
		switch fnType.In(i).Kind() {
		case reflect.Float64:
			paramTypes = append(paramTypes, wasm.ValueTypeF64)
		case reflect.Float32:
			paramTypes = append(paramTypes, wasm.ValueTypeF32)
		case reflect.Int64:
			paramTypes = append(paramTypes, wasm.ValueTypeI64)
		default:
			paramTypes = append(paramTypes, wasm.ValueTypeI32)
		}
	}

	resultTypes := make([]wasm.ValueType, 0, numResults)
	for i := 0; i < numResults; i++ {
		if hasErrorResult && i == numResults-1 {
			continue
		}
		switch fnType.Out(i).Kind() {
		case reflect.Float64:
			resultTypes = append(resultTypes, wasm.ValueTypeF64)
		case reflect.Float32:
			resultTypes = append(resultTypes, wasm.ValueTypeF32)
		case reflect.Int64:
			resultTypes = append(resultTypes, wasm.ValueTypeI64)
		default:
			resultTypes = append(resultTypes, wasm.ValueTypeI32)
		}
	}

	// Create the host function object
	hf := &HostFunction{
		module:          modName,
		name:            funcName,
		wasmParamTypes:  paramTypes,
		wasmResultTypes: resultTypes,
		messages:        make(map[hfMessageType]string),
	}
	for _, opt := range opts {
		opt(hf)
	}

	// Make the host function wrapper
	hf.function = wasm.GoFunc(func(ctx context.Context, stack []uint64) {

		// Log any panics that occur in the host function
		defer func() {
			if e := utils.ConvertToError(recover()); e != nil {
				logger.Err(ctx, e).Msg("Panic in host function.")
			}
		}()

		// Get the host function's metadata
		fullName := modName + "." + funcName
		fnMeta := metadata.GetFunctionImportMetadata(ctx, fullName)

		// Read input parameter values
		params := make([]any, 0, numParams)
		for i := 0; i < numParams; i++ {
			if hasContextParam && i == 0 {
				continue
			}
			paramType := fnType.In(i)
			var rvParam reflect.Value
			if paramType.Kind() == reflect.Ptr {
				rvParam = reflect.New(paramType.Elem())
			} else {
				rvParam = reflect.New(paramType)
			}
			params = append(params, rvParam.Interface())
		}
		if err := decodeParams(ctx, fnMeta, stack, params...); err != nil {
			logger.Err(ctx, err).Msg("Error decoding input parameters.")
			return
		}

		// Prepare to call the host function
		results := make([]any, 0, numResults)
		wrappedFn := func() error {

			// prepare the input parameters
			in := make([]reflect.Value, 0, numParams)
			i := 0
			if hasContextParam {
				in = append(in, reflect.ValueOf(ctx))
				i++
			}
			for _, param := range params {
				paramType := fnType.In(i)
				if paramType.Kind() == reflect.Ptr {
					in = append(in, reflect.ValueOf(param))
				} else {
					in = append(in, reflect.ValueOf(param).Elem())
				}
				i++
			}

			// invoke the function
			out := rvFunc.Call(in)

			// check for an error
			if hasErrorResult && len(out) > 0 {
				if err, ok := out[len(out)-1].Interface().(error); ok && err != nil {
					return err
				}
			}

			// copy results to the results slice
			for i := 0; i < numResults; i++ {
				if hasErrorResult && i == numResults-1 {
					continue
				} else {
					results = append(results, out[i].Interface())
				}
			}

			return nil
		}

		// Call the host function
		if ok := callHostFunction(ctx, wrappedFn, hf.messages); !ok {
			return
		}

		// Encode the results (if there are any)
		if len(results) > 0 {
			if err := encodeResults(ctx, fnMeta, stack, results...); err != nil {
				logger.Err(ctx, err).Msg("Error encoding results.")
			}
		}
	})

	return hf, nil
}

func (host *WasmHost) instantiateHostFunctions(ctx context.Context) error {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	hostFnsByModule := make(map[string][]*HostFunction)
	for _, hf := range host.hostFunctions {
		hostFnsByModule[hf.module] = append(hostFnsByModule[hf.module], hf)
	}

	for module, modHostFns := range hostFnsByModule {
		builder := host.runtime.NewHostModuleBuilder(module)
		for _, hf := range modHostFns {
			builder.NewFunctionBuilder().
				WithGoFunction(hf.function, hf.wasmParamTypes, hf.wasmResultTypes).
				WithName(module + "." + hf.name).
				Export(hf.name)
		}
		if _, err := builder.Instantiate(ctx); err != nil {
			return fmt.Errorf("failed to instantiate the %s host module: %w", module, err)
		}
	}

	return nil
}

func decodeParams(ctx context.Context, fn *metadata.Function, stack []uint64, params ...any) error {

	// regardless of the outcome, ensure parameter values are cleared from the stack before returning
	indirect := false
	defer func() {
		i := 0
		if indirect {
			i = 1
		}
		for ; i < len(stack); i++ {
			stack[i] = 0
		}
	}()

	var expected int
	if fn != nil {
		expected = len(fn.Parameters)
	} else {
		expected = len(stack)
	}
	if len(params) != expected {
		return fmt.Errorf("expected %d parameters, but got %d", expected, len(params))
	}

	wa, err := languages.GetWasmAdapter(ctx)
	if err != nil {
		return err
	}

	getEncLen := func(i int) (int, error) {
		if fn == nil {
			return 1, nil
		}

		typ := fn.Parameters[i].Type
		encLength, err := wa.GetEncodingLength(ctx, typ)
		if err != nil {
			return 0, err
		}

		return encLength, nil
	}

	totalEncLength := 0
	encLengths := make([]int, len(params))
	for i := 0; i < len(params); i++ {
		if n, err := getEncLen(i); err != nil {
			return err
		} else {
			encLengths[i] = n
			totalEncLength += n
		}
	}

	stackPos := 0

	// handle result indirection, if supported and needed
	if _, ok := wa.(languages.WasmAdapterWithIndirection); ok {
		if len(stack) == totalEncLength+1 {
			stackPos++
			indirect = true
		}
	}

	for i, p := range params {
		// get values from the stack
		encLength := encLengths[i]
		vals := stack[stackPos : stackPos+encLength]
		stackPos += encLength

		// decode the values for the parameter
		var typ string
		if fn != nil {
			typ = fn.Parameters[i].Type
		}
		if err := wa.DecodeData(ctx, typ, vals, p); err != nil {
			return err
		}
	}

	return nil
}

func encodeResults(ctx context.Context, fn *metadata.Function, stack []uint64, results ...any) error {

	if fn != nil {
		expected := len(fn.Results)
		if len(results) != expected {
			return fmt.Errorf("expected %d results, but got %d", expected, len(results))
		}
	}

	wa, err := languages.GetWasmAdapter(ctx)
	if err != nil {
		return err
	}

	// if result indirection is supported, write the results indirectly if warranted
	if wa, ok := wa.(languages.WasmAdapterWithIndirection); ok {
		if len(stack) > 0 && stack[0] != 0 {
			return wa.WriteIndirectResults(ctx, fn, uint32(stack[0]), results)
		}
	}

	cleaner := utils.NewCleanerN(len(results))
	errs := make([]error, 0, len(results))

	stackPos := 0
	for i, r := range results {

		// get the type of the result
		var typ string
		if fn != nil {
			typ = fn.Results[i].Type
		}

		// encode the result
		vals, cln, err := wa.EncodeData(ctx, typ, r)
		cleaner.AddCleaner(cln)
		if err != nil {
			errs = append(errs, err)
		}

		// push the encoded values onto the stack
		for _, v := range vals {
			stack[stackPos] = v
			stackPos++
		}
	}

	// unpin any values pinned during encoding
	if err := cleaner.Clean(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func callHostFunction(ctx context.Context, fn func() error, msgs map[hfMessageType]string) bool {

	detail := msgs[hfMessageDetail]

	if msg, ok := msgs[hfMessageStarting]; ok {
		l := logger.Info(ctx).Bool("user_visible", true)
		if detail != "" {
			l.Str("detail", detail)
		}
		l.Msg(msg)
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	if errors.Is(err, context.Canceled) {
		if msg, ok := msgs[hfMessageCancelled]; ok {
			l := logger.Warn(ctx).Bool("user_visible", true).Dur("duration_ms", duration)
			if detail != "" {
				l.Str("detail", detail)
			}
			l.Msg(msg)
		}
		return false
	} else if err != nil {
		if msg, ok := msgs[hfMessageError]; ok {
			l := logger.Err(ctx, err).Bool("user_visible", true).Dur("duration_ms", duration)
			if detail != "" {
				l.Str("detail", detail)
			}
			l.Msg(msg)
		}
		return false
	} else {
		if msg, ok := msgs[hfMessageCompleted]; ok {
			l := logger.Info(ctx).Bool("user_visible", true).Dur("duration_ms", duration)
			if detail != "" {
				l.Str("detail", detail)
			}
			l.Msg(msg)
		}
		return true
	}
}
