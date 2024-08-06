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

	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func (wa *wasmAdapter) readObject(ctx context.Context, mem wasm.Memory, typeName string, offset uint32) (data any, err error) {
	switch typeName {
	case "~lib/arraybuffer/ArrayBuffer":
		return wa.readBytes(mem, offset)
	case "string", "~lib/string/String":
		return wa.readString(mem, offset)
	case "~lib/date/Date", "~lib/wasi_date/wasi_Date":
		return wa.readDate(mem, offset)
	}

	typ, err := wa.typeInfo.getTypeDefinition(ctx, typeName)
	if err != nil {
		return nil, err
	}

	id, _ := mem.ReadUint32Le(offset - 8)
	if id != typ.Id {
		return nil, fmt.Errorf("pointer is not to a %s", typeName)
	}

	if wa.typeInfo.IsArrayType(typeName) {
		return wa.readArray(ctx, mem, typ, offset)
	} else if wa.typeInfo.IsMapType(typeName) {
		return wa.readMap(ctx, mem, typ, offset)
	}

	return wa.readClass(ctx, mem, typ, offset)
}

func (wa *wasmAdapter) writeObject(ctx context.Context, mod wasm.Module, typeName string, val any) (offset uint32, err error) {
	switch typeName {
	case "~lib/arraybuffer/ArrayBuffer":
		switch v := val.(type) {
		case []byte:
			return wa.writeBytes(ctx, mod, v)
		case *[]byte:
			if v == nil {
				return 0, nil
			}
			return wa.writeBytes(ctx, mod, *v)
		default:
			return 0, fmt.Errorf("input value is not a byte array")
		}

	case "string", "~lib/string/String":
		switch v := val.(type) {
		case string:
			return wa.writeString(ctx, mod, v)
		case *string:
			if v == nil {
				return 0, nil
			}
			return wa.writeString(ctx, mod, *v)
		default:
			return 0, fmt.Errorf("input value is not a string")
		}

	case "~lib/date/Date", "~lib/wasi_date/wasi_Date":
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

		return wa.writeDate(ctx, mod, t)
	}

	typ, err := wa.typeInfo.getTypeDefinition(ctx, typeName)
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

	if wa.typeInfo.IsArrayType(typeName) {
		return wa.writeArray(ctx, mod, typ, val)
	} else if wa.typeInfo.IsMapType(typeName) {
		return wa.writeMap(ctx, mod, typ, val)
	} else {
		return wa.writeClass(ctx, mod, typ, val)
	}
}
