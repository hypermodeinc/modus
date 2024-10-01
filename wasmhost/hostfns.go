/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"time"

	"hypruntime/langsupport"
	"hypruntime/logger"
	"hypruntime/plugins"
	"hypruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

var rtContext = reflect.TypeFor[context.Context]()
var rtError = reflect.TypeFor[error]()

type hfMessages struct {
	msgStarting  string
	msgCompleted string
	msgCancelled string
	msgError     string

	fnDetail  any
	msgDetail string
}

type hostFunction struct {
	module          string
	name            string
	function        wasm.GoFunction
	wasmParamTypes  []wasm.ValueType
	wasmResultTypes []wasm.ValueType
	messages        *hfMessages
}

func (hf *hostFunction) Name() string {
	return hf.module + "." + hf.name
}

type HostFunctionOption func(*hostFunction)

func WithStartingMessage(text string) HostFunctionOption {
	return func(hf *hostFunction) {
		hf.messages.msgStarting = text
	}
}

func WithCompletedMessage(text string) HostFunctionOption {
	return func(hf *hostFunction) {
		hf.messages.msgCompleted = text
	}
}

func WithCancelledMessage(text string) HostFunctionOption {
	return func(hf *hostFunction) {
		hf.messages.msgCancelled = text
	}
}

func WithErrorMessage(text string) HostFunctionOption {
	return func(hf *hostFunction) {
		hf.messages.msgError = text
	}
}

func WithMessageDetail(fn any) HostFunctionOption {
	return func(hf *hostFunction) {
		hf.messages.fnDetail = fn
	}
}

func (host *wasmHost) RegisterHostFunction(modName, funcName string, fn any, opts ...HostFunctionOption) error {
	hf, err := host.newHostFunction(modName, funcName, fn, opts...)
	if err != nil {
		return fmt.Errorf("failed to register host function %s.%s: %w", modName, funcName, err)
	}

	host.hostFunctions = append(host.hostFunctions, hf)
	return nil
}

func (host *wasmHost) newHostFunction(modName, funcName string, fn any, opts ...HostFunctionOption) (*hostFunction, error) {
	fullName := modName + "." + funcName
	rvFunc := reflect.ValueOf(fn)
	if rvFunc.Kind() != reflect.Func {
		return nil, fmt.Errorf("host function %s is not a function type", fullName)
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
		return nil, fmt.Errorf("host function %s cannot return multiple data values", fullName)
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
	hf := &hostFunction{
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
			return nil, fmt.Errorf("message detail func for host function %s is not a function type", fullName)
		}
		rtDetail = rvDetail.Type()
		if rtDetail.NumOut() != 1 || rtDetail.Out(0).Kind() != reflect.String {
			return nil, fmt.Errorf("message detail func for host function %s must have one string return value", fullName)
		}

		start, end := 0, rtDetail.NumIn()
		if hasContextParam {
			start++
			end++
		}
		i := 0
		for j := start; j < end; j++ {
			if rtDetail.In(i) != rtFunc.In(j) {
				return nil, fmt.Errorf("message detail func for host function %s has mismatched parameter types", fullName)
			}
			i++
		}
	}

	// Make the host function wrapper
	hf.function = wasm.GoFunc(func(ctx context.Context, stack []uint64) {
		span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
		defer span.Finish()

		// Log any panics that occur in the host function
		defer func() {
			if r := recover(); r != nil {
				logger.Err(ctx, utils.ConvertToError(r)).Str("host_function", fullName).Msg("Panic in host function.")
				if utils.HypermodeDebugEnabled() {
					debug.PrintStack()
				}
			}
		}()

		// Get the plugin of the function that invoked this host function
		plugin, ok := plugins.GetPluginFromContext(ctx)
		if !ok {
			logger.Error(ctx).Str("host_function", fullName).Msg("Plugin not found in context.")
			return
		}

		// Get the execution plan for the host function
		plan, ok := plugin.ExecutionPlans[fullName]
		if !ok {
			logger.Error(ctx).Str("host_function", fullName).Msg("Execution plan not found.")
		}

		// Get the Wasm adapter
		wa, err := langsupport.GetWasmAdapter(ctx)
		if err != nil {
			logger.Err(ctx, err).Msg("Error getting Wasm adapter.")
			return
		}

		// Read input parameter values
		params := make([]any, 0, numParams)
		for i := 0; i < numParams; i++ {
			if hasContextParam && i == 0 {
				continue
			}
			rtParam := rtFunc.In(i)
			rvParam := reflect.New(rtParam).Elem()
			params = append(params, rvParam.Interface())
		}
		if err := decodeParams(ctx, wa, plan, stack, params); err != nil {
			logger.Err(ctx, err).Str("host_function", fullName).Any("data", params).Msg("Error decoding input parameters.")
			return
		}

		// prepare the input parameters
		inputs := make([]reflect.Value, 0, numParams)
		if hasContextParam {
			inputs = append(inputs, reflect.ValueOf(ctx))
		}
		for _, param := range params {
			inputs = append(inputs, reflect.ValueOf(param))
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
			if err := encodeResults(ctx, wa, plan, stack, results); err != nil {
				logger.Err(ctx, err).Str("host_function", fullName).Any("data", results).Msg("Error encoding results.")
			}
		}
	})

	return hf, nil
}

func (host *wasmHost) instantiateHostFunctions(ctx context.Context) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	hostFnsByModule := make(map[string][]*hostFunction)
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

func decodeParams(ctx context.Context, wa langsupport.WasmAdapter, plan langsupport.ExecutionPlan, stack []uint64, params []any) error {

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

	expected := len(plan.ParamHandlers())
	if len(params) != expected {
		return fmt.Errorf("expected %d parameters, but got %d", expected, len(params))
	}

	stackPos := 0
	if plan.UseResultIndirection() {
		indirect = true
		stackPos++
	}

	for i, handler := range plan.ParamHandlers() {
		encLength := int(handler.TypeInfo().EncodingLength())
		vals := stack[stackPos : stackPos+encLength]
		stackPos += encLength

		data, err := handler.Decode(ctx, wa, vals)
		if err != nil {
			return err
		}

		// special case for structs represented as maps
		switch m := data.(type) {
		case map[string]any:
			if _, ok := (params[i]).(map[string]any); !ok {
				if err := utils.MapToStruct(m, &params[i]); err != nil {
					return err
				}
				continue
			}
		case *map[string]any:
			if _, ok := (params[i]).(*map[string]any); !ok {
				if err := utils.MapToStruct(*m, &params[i]); err != nil {
					return err
				}
				continue
			}
		}

		// special case for pointers that need to be dereferenced
		if handler.TypeInfo().ReflectedType().Kind() == reflect.Ptr && reflect.TypeOf(params[i]).Kind() != reflect.Ptr {
			params[i] = utils.DereferencePointer(data)
			continue
		}

		// special case for non-pointers that need to be converted to pointers
		if handler.TypeInfo().ReflectedType().Kind() != reflect.Ptr && reflect.TypeOf(params[i]).Kind() == reflect.Ptr {
			params[i] = utils.MakePointer(data)
			continue
		}

		params[i] = data
	}

	return nil
}

func encodeResults(ctx context.Context, wa langsupport.WasmAdapter, plan langsupport.ExecutionPlan, stack []uint64, results []any) error {

	expected := len(plan.ResultHandlers())
	if len(results) != expected {
		return fmt.Errorf("expected %d results, but got %d", expected, len(results))
	}

	if plan.UseResultIndirection() {
		return writeIndirectResults(ctx, wa, plan, uint32(stack[0]), results)
	}

	cleaner := utils.NewCleanerN(len(results))
	stackPos := 0

	for i, handler := range plan.ResultHandlers() {
		vals, cln, err := handler.Encode(ctx, wa, results[i])
		cleaner.AddCleaner(cln)
		if err != nil {
			return err
		}

		for _, v := range vals {
			stack[stackPos] = v
			stackPos++
		}
	}

	return cleaner.Clean()
}

func writeIndirectResults(ctx context.Context, wa langsupport.WasmAdapter, plan langsupport.ExecutionPlan, offset uint32, results []any) (err error) {

	// multiple-results are written like a struct

	cleaner := utils.NewCleanerN(len(results))
	defer func() {
		if e := cleaner.Clean(); e != nil && err == nil {
			err = e
		}
	}()

	handlers := plan.ResultHandlers()

	fieldOffset := uint32(0)
	for i, handler := range handlers {
		size := handler.TypeInfo().Size()
		fieldType := handler.TypeInfo().Name()
		alignment, err := wa.TypeInfo().GetAlignmentOfType(ctx, fieldType)
		if err != nil {
			return err
		}

		fieldOffset = langsupport.AlignOffset(fieldOffset, alignment)

		cln, err := handler.Write(ctx, wa, offset+fieldOffset, results[i])
		cleaner.AddCleaner(cln)
		if err != nil {
			return err
		}

		fieldOffset += size
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
