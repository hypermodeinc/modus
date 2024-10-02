/*
 * Copyright 2024 Hypermode, Inc.
 */

package datasource

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	"github.com/buger/jsonparser"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/datasource/httpclient"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

const DataSourceName = "HypermodeFunctionsDataSource"

type callInfo struct {
	Function   fieldInfo      `json:"fn"`
	Parameters map[string]any `json:"data"`
}

type HypermodeDataSource struct {
	WasmHost wasmhost.WasmHost
}

func (ds *HypermodeDataSource) Load(ctx context.Context, input []byte, out *bytes.Buffer) error {

	// Parse the input to get the function call info
	var ci callInfo
	err := utils.JsonDeserialize(input, &ci)
	if err != nil {
		return fmt.Errorf("error parsing input: %w", err)
	}

	// Load the data
	result, gqlErrors, err := ds.callFunction(ctx, &ci)

	// Write the response
	err = writeGraphQLResponse(ctx, out, result, gqlErrors, err, &ci)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Error creating GraphQL response.")
	}

	return err
}

func (*HypermodeDataSource) LoadWithFiles(ctx context.Context, input []byte, files []httpclient.File, out *bytes.Buffer) (err error) {
	// See https://github.com/wundergraph/graphql-go-tools/pull/758
	panic("not implemented")
}

func (ds *HypermodeDataSource) callFunction(ctx context.Context, callInfo *callInfo) (any, []resolve.GraphQLError, error) {

	// Get the function info
	fnInfo, err := ds.WasmHost.GetFunctionInfo(callInfo.Function.Name)
	if err != nil {
		return nil, nil, err
	}

	// Call the function
	execInfo, err := ds.WasmHost.CallFunction(ctx, fnInfo, callInfo.Parameters)
	if err != nil {
		// The full error message has already been logged.  Return a generic error to the caller, which will be included in the response.
		return nil, nil, errors.New("error calling function")
	}

	// Store the execution info into the function output map.
	outputMap := ctx.Value(utils.FunctionOutputContextKey).(map[string]wasmhost.ExecutionInfo)
	outputMap[callInfo.Function.AliasOrName()] = execInfo

	// Transform messages (and error lines in the output buffers) to GraphQL errors.
	messages := append(execInfo.Messages(), utils.TransformConsoleOutput(execInfo.Buffers())...)
	gqlErrors := transformErrors(messages, callInfo)

	// Get the result.
	result := execInfo.Result()

	// If we have multiple results, unpack them into a map that matches the schema generated type.
	if results, ok := result.([]any); ok && len(fnInfo.ExecutionPlan().ResultHandlers()) > 1 {
		fnMeta := fnInfo.Metadata()
		m := make(map[string]any, len(results))
		for i, r := range results {
			name := fnMeta.Results[i].Name
			if name == "" {
				name = fmt.Sprintf("item%d", i+1)
			}
			m[name] = r
		}
		result = m
	}

	return result, gqlErrors, err
}

func writeGraphQLResponse(ctx context.Context, out *bytes.Buffer, result any, gqlErrors []resolve.GraphQLError, fnErr error, ci *callInfo) error {

	fieldName := ci.Function.AliasOrName()

	// Include the function error
	if fnErr != nil {
		gqlErrors = append(gqlErrors, resolve.GraphQLError{
			Message: fnErr.Error(),
			Path:    []any{fieldName},
			Extensions: map[string]interface{}{
				"level": "error",
			},
		})
	}

	// If there are GraphQL errors, serialize them as json
	var jsonErrors []byte
	if len(gqlErrors) > 0 {
		var err error
		jsonErrors, err = utils.JsonSerialize(gqlErrors)
		if err != nil {
			return err
		}
	}

	// If there is any result data, or if the data is null without errors, serialize the data as json
	var jsonData []byte
	if result != nil || len(gqlErrors) == 0 {
		jsonResult, err := utils.JsonSerialize(result)
		if err != nil {
			if err, ok := err.(*json.UnsupportedValueError); ok {
				msg := fmt.Sprintf("Function completed successfully, but the result contains a %v value that cannot be serialized to JSON.", err.Value)
				logger.Warn(ctx).
					Bool("user_visible", true).
					Str("function", ci.Function.Name).
					Str("result", fmt.Sprintf("%+v", result)).
					Msg(msg)
				fmt.Fprintf(out, `{"errors":[{"message":"%s","path":["%s"],"extensions":{"level":"error"}}]}`, msg, fieldName)
				return nil
			}
			return err
		}

		// Transform the data
		if r, err := transformValue(jsonResult, &ci.Function); err != nil {
			return err
		} else {
			jsonData = r
		}
	}

	// Write the response.  This should be as efficient as possible, as it is called for every function invocation.
	out.Grow(len(jsonData) + len(jsonErrors) + len(fieldName) + 26)
	out.WriteByte('{')
	if len(jsonData) > 0 {
		out.WriteString(`"data":{"`)
		out.WriteString(fieldName)
		out.WriteString(`":`)
		out.Write(jsonData)
		out.WriteByte('}')
	}
	if len(jsonErrors) > 0 {
		if len(jsonData) > 0 {
			out.WriteByte(',')
		}
		out.WriteString(`"errors":`)
		out.Write(jsonErrors)
	}
	out.WriteByte('}')

	return nil
}

var nullWord = []byte("null")

func transformValue(data []byte, tf *fieldInfo) (result []byte, err error) {
	if len(tf.Fields) == 0 || len(data) == 0 || bytes.Equal(data, nullWord) {
		return data, nil
	}

	switch data[0] {
	case '{':
		if tf.IsMapType {
			return transformMap(data, tf)
		} else {
			return transformObject(data, tf)
		}
	case '[':
		return transformArray(data, tf)
	default:
		return nil, fmt.Errorf("expected object or array")
	}
}

func transformArray(data []byte, tf *fieldInfo) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('[')

	var loopErr error
	_, err := jsonparser.ArrayEach(data, func(val []byte, _ jsonparser.ValueType, _ int, _ error) {
		if loopErr != nil {
			return
		}
		val, err := transformValue(val, tf)
		if err != nil {
			loopErr = err
			return
		}
		if buf.Len() > 1 {
			buf.WriteByte(',')
		}
		buf.Write(val)
	})
	if err != nil {
		return nil, err
	}
	if loopErr != nil {
		return nil, loopErr
	}

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func transformObject(data []byte, tf *fieldInfo) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('{')
	for i, f := range tf.Fields {
		var val []byte
		if f.Name == "__typename" {
			val = []byte(`"` + tf.TypeName + `"`)
		} else {
			v, dataType, _, err := jsonparser.Get(data, f.Name)
			if err != nil {
				return nil, err
			}
			if dataType == jsonparser.String {
				// Note, string values here will be escaped for internal quotes, newlines, etc.,
				// but will be missing outer quotes.  So we need to add them back.
				v = []byte(`"` + string(v) + `"`)
			}
			val, err = transformValue(v, &f)
			if err != nil {
				return nil, err
			}
		}
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('"')
		buf.WriteString(f.AliasOrName())
		buf.WriteString(`":`)
		buf.Write(val)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func transformMap(data []byte, tf *fieldInfo) ([]byte, error) {

	// check for pseudo map
	md, dt, _, err := jsonparser.Get(data, "$mapdata")
	if err == nil && dt == jsonparser.Array {
		return transformPseudoMap(md, tf)
	}

	var keyType string
	for _, f := range tf.Fields {
		if f.Name == "key" {
			keyType = f.TypeName
			break
		}
	}

	buf := bytes.Buffer{}
	buf.WriteByte('[')
	if err := jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if buf.Len() > 1 {
			buf.WriteByte(',')
		}

		b := bytes.Buffer{}
		b.WriteByte('{')
		b.WriteString(`"key":`)
		if keyType == "String" {
			k, err := utils.JsonSerialize(string(key))
			if err != nil {
				return err
			}
			b.Write(k)
		} else {
			b.Write(key)
		}
		b.WriteString(`,"value":`)
		if dataType == jsonparser.String {
			b.WriteString(`"`)
			b.Write(value)
			b.WriteString(`"`)
		} else {
			b.Write(value)
		}
		b.WriteByte('}')

		val, err := transformObject(b.Bytes(), tf)
		if err != nil {
			return err
		}
		buf.Write(val)

		return nil
	}); err != nil {
		return nil, err
	}

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func transformPseudoMap(data []byte, tf *fieldInfo) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('[')

	var loopErr error
	_, err := jsonparser.ArrayEach(data, func(item []byte, _ jsonparser.ValueType, _ int, _ error) {
		if loopErr != nil {
			return
		}

		key, kdt, _, err := jsonparser.Get(item, "key")
		if err != nil {
			loopErr = err
			return
		}

		value, vdt, _, err := jsonparser.Get(item, "value")
		if err != nil {
			loopErr = err
			return
		}

		if buf.Len() > 1 {
			buf.WriteByte(',')
		}

		b := bytes.Buffer{}
		b.WriteByte('{')
		b.WriteString(`"key":`)
		if kdt == jsonparser.String {
			b.WriteString(`"`)
			b.Write(key)
			b.WriteString(`"`)
		} else {
			b.Write(key)
		}
		b.WriteString(`,"value":`)
		if vdt == jsonparser.String {
			b.WriteString(`"`)
			b.Write(value)
			b.WriteString(`"`)
		} else {
			b.Write(value)
		}
		b.WriteByte('}')

		val, err := transformObject(b.Bytes(), tf)
		if err != nil {
			loopErr = err
			return
		}
		buf.Write(val)
	})
	if err != nil {
		return nil, err
	}
	if loopErr != nil {
		return nil, loopErr
	}

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func transformErrors(messages []utils.LogMessage, ci *callInfo) []resolve.GraphQLError {
	errors := make([]resolve.GraphQLError, 0, len(messages))
	for _, msg := range messages {
		// Only include errors.  Other messages will be captured later and
		// passed back as logs in the extensions section of the response.
		if msg.IsError() {
			errors = append(errors, resolve.GraphQLError{
				Message: msg.Message,
				Path:    []any{ci.Function.AliasOrName()},
				Extensions: map[string]interface{}{
					"level": msg.Level,
				},
			})
		}
	}
	return errors
}
