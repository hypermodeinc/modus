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

var rtContext = reflect.TypeOf((*context.Context)(nil)).Elem()
var rtError = reflect.TypeOf((*error)(nil)).Elem()

type hfMessages struct {
	msgStarting  string
	msgCompleted string
	msgCancelled string
	msgError     string

	fnDetail  any
	msgDetail string
}

type HostFunction struct {
	module          string
	name            string
	function        wasm.GoFunction
	wasmParamTypes  []wasm.ValueType
	wasmResultTypes []wasm.ValueType
	messages        *hfMessages
}

func (hf *HostFunction) Name() string {
	return hf.module + "." + hf.name
}

func WithStartingMessage(text string) func(*HostFunction) {
	return func(hf *HostFunction) {
		hf.messages.msgStarting = text
	}
}

func WithCompletedMessage(text string) func(*HostFunction) {
	return func(hf *HostFunction) {
		hf.messages.msgCompleted = text
	}
}

func WithCancelledMessage(text string) func(*HostFunction) {
	return func(hf *HostFunction) {
		hf.messages.msgCancelled = text
	}
}

func WithErrorMessage(text string) func(*HostFunction) {
	return func(hf *HostFunction) {
		hf.messages.msgError = text
	}
}

func WithMessageDetail(fn any) func(*HostFunction) {
	return func(hf *HostFunction) {
		hf.messages.fnDetail = fn
	}
}

func (host *WasmHost) RegisterHostFunction(modName, funcName string, fn any, opts ...func(*HostFunction)) error {
	hf, err := newHostFunction(modName, funcName, fn, opts...)
	if err != nil {
		return fmt.Errorf("failed to register host function %s.%s: %w", modName, funcName, err)
	}

	host.hostFunctions = append(host.hostFunctions, hf)
	return nil
}

func newHostFunction(modName, funcName string, fn any, opts ...func(*HostFunction)) (*HostFunction, error) {
	rvFunc := reflect.ValueOf(fn)
	if rvFunc.Kind() != reflect.Func {
		return nil, fmt.Errorf("host function %s.%s is not a function type", modName, funcName)
	}

	rtFunc := rvFunc.Type()
	numParams := rtFunc.NumIn()
	numResults := rtFunc.NumOut()

	// Optionally, the first parameter can be a context.
	var hasContextParam bool
	if numParams > 0 && rtFunc.In(0).Implements(rtContext) {
		hasContextParam = true
	}

	// Optionally, the last return value can be an error.
	var hasErrorResult bool
	if numResults > 0 && rtFunc.Out(numResults-1).Implements(rtError) {
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
		switch rtFunc.In(i).Kind() {
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
		switch rtFunc.Out(i).Kind() {
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
		messages:        &hfMessages{},
	}
	for _, opt := range opts {
		opt(hf)
	}

	// Prep the message detail function
	var rvDetail reflect.Value
	var rtDetail reflect.Type
	if hf.messages.fnDetail != nil {
		rvDetail = reflect.ValueOf(hf.messages.fnDetail)
		if rvDetail.Kind() != reflect.Func {
			return nil, fmt.Errorf("message detail func for host function %s.%s is not a function type", modName, funcName)
		}
		rtDetail = rvDetail.Type()
		if rtDetail.NumOut() != 1 || rtDetail.Out(0).Kind() != reflect.String {
			return nil, fmt.Errorf("message detail func for host function %s.%s must have one string return value", modName, funcName)
		}

		start, end := 0, rtDetail.NumIn()
		if hasContextParam {
			start++
			end++
		}
		i := 0
		for j := start; j < end; j++ {
			if rtDetail.In(i) != rtFunc.In(j) {
				return nil, fmt.Errorf("message detail func for host function %s.%s has mismatched parameter types", modName, funcName)
			}
			i++
		}
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
			rtParam := rtFunc.In(i)
			var rvParam reflect.Value
			if rtParam.Kind() == reflect.Ptr {
				rvParam = reflect.New(rtParam.Elem())
			} else {
				rvParam = reflect.New(rtParam)
			}
			params = append(params, rvParam.Interface())
		}
		if err := decodeParams(ctx, fnMeta, stack, params...); err != nil {
			logger.Err(ctx, err).Msg("Error decoding input parameters.")
			return
		}

		// prepare the input parameters
		inputs := make([]reflect.Value, 0, numParams)
		i := 0
		if hasContextParam {
			inputs = append(inputs, reflect.ValueOf(ctx))
			i++
		}
		for _, param := range params {
			rtParam := rtFunc.In(i)
			if rtParam.Kind() == reflect.Ptr {
				inputs = append(inputs, reflect.ValueOf(param))
			} else {
				inputs = append(inputs, reflect.ValueOf(param).Elem())
			}
			i++
		}

		// Prepare to call the host function
		results := make([]any, 0, numResults)
		wrappedFn := func() error {

			// invoke the function
			out := rvFunc.Call(inputs)

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

		// If there is a message detail function, call it to get the detail message
		msgs := *hf.messages
		if msgs.fnDetail != nil {
			start, end := 0, rtDetail.NumIn()
			if hasContextParam {
				start++
				end++
			}
			msgs.msgDetail = rvDetail.Call(inputs[start:end])[0].String()
		}

		// Call the host function
		if ok := callHostFunction(ctx, wrappedFn, msgs); !ok {
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

func callHostFunction(ctx context.Context, fn func() error, msgs hfMessages) bool {
	if msgs.msgStarting != "" {
		l := logger.Info(ctx).Bool("user_visible", true)
		if msgs.msgDetail != "" {
			l.Str("detail", msgs.msgDetail)
		}
		l.Msg(msgs.msgStarting)
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	if errors.Is(err, context.Canceled) {
		if msgs.msgCancelled != "" {
			l := logger.Warn(ctx).Bool("user_visible", true).Dur("duration_ms", duration)
			if msgs.msgDetail != "" {
				l.Str("detail", msgs.msgDetail)
			}
			l.Msg(msgs.msgCancelled)
		}
		return false
	} else if err != nil {
		if msgs.msgError != "" {
			l := logger.Err(ctx, err).Bool("user_visible", true).Dur("duration_ms", duration)
			if msgs.msgDetail != "" {
				l.Str("detail", msgs.msgDetail)
			}
			l.Msg(msgs.msgError)
		}
		return false
	} else {
		if msgs.msgCompleted != "" {
			l := logger.Info(ctx).Bool("user_visible", true).Dur("duration_ms", duration)
			if msgs.msgDetail != "" {
				l.Str("detail", msgs.msgDetail)
			}
			l.Msg(msgs.msgCompleted)
		}
		return true
	}
}
