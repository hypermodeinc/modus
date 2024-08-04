/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"hmruntime/plugins"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func readObject(ctx context.Context, mem wasm.Memory, typ plugins.TypeInfo, offset uint32) (data any, err error) {
	switch typ.Name {
	case "ArrayBuffer":
		return readBytes(mem, offset)
	case "string":
		return readString(mem, offset)
	case "Date":
		return readDate(mem, offset)
	}

	def, err := getTypeDefinition(ctx, typ.Path)
	if err != nil {
		return nil, err
	}

	id, _ := mem.ReadUint32Le(offset - 8)
	if id != def.Id {
		return nil, fmt.Errorf("pointer is not to a %s", typ.Name)
	}

	if isArrayType(typ.Path) {
		return readArray(ctx, mem, def, offset)
	} else if isMapType(typ.Path) {
		return readMap(ctx, mem, def, offset)
	}

	return readClass(ctx, mem, def, offset)
}

func writeObject(ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, val any) (offset uint32, err error) {
	switch typ.Name {
	case "ArrayBuffer":
		switch v := val.(type) {
		case []byte:
			return writeBytes(ctx, mod, v)
		case *[]byte:
			if v == nil {
				return 0, nil
			}
			return writeBytes(ctx, mod, *v)
		default:
			return 0, fmt.Errorf("input value is not a byte array")
		}

	case "string":
		switch v := val.(type) {
		case string:
			return writeString(ctx, mod, v)
		case *string:
			if v == nil {
				return 0, nil
			}
			return writeString(ctx, mod, *v)
		default:
			return 0, fmt.Errorf("input value is not a string")
		}

	case "Date":
		var t time.Time
		switch v := val.(type) {
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return 0, err
			}
			t = time.UnixMilli(n)
		case *json.Number:
			if v == nil {
				return 0, nil
			}
			n, err := v.Int64()
			if err != nil {
				return 0, err
			}
			t = time.UnixMilli(n)
		case string:
			var err error
			t, err = utils.ParseTime(v)
			if err != nil {
				return 0, err
			}
		case *string:
			if v == nil {
				return 0, nil
			}
			var err error
			t, err = utils.ParseTime(*v)
			if err != nil {
				return 0, err
			}
		case utils.JSONTime:
			t = time.Time(v)
		case *utils.JSONTime:
			if v == nil {
				return 0, nil
			}
			t = time.Time(*v)
		case time.Time:
			t = v
		case *time.Time:
			if v == nil {
				return 0, nil
			}
			t = *v
		default:
			return 0, fmt.Errorf("input value is not a valid for a time object")
		}

		return writeDate(ctx, mod, t)
	}

	def, err := getTypeDefinition(ctx, typ.Path)
	if err != nil {
		return 0, err
	}

	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		v := reflect.ValueOf(val)
		if v.IsNil() {
			return 0, nil
		} else {
			val = v.Elem().Interface()
		}
	}

	if isArrayType(typ.Path) {
		return writeArray(ctx, mod, def, val)
	} else if isMapType(typ.Path) {
		return writeMap(ctx, mod, def, val)
	} else {
		return writeClass(ctx, mod, def, val)
	}
}
