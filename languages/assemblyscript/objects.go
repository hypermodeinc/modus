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
)

func (wa *wasmAdapter) readObject(ctx context.Context, typ string, offset uint32) (data any, err error) {
	if offset == 0 {
		return nil, nil
	}

	switch typ {
	case "~lib/arraybuffer/ArrayBuffer":
		return wa.readBytes(offset)
	case "string", "~lib/string/String":
		return wa.readString(offset)
	case "~lib/date/Date", "~lib/wasi_date/wasi_Date":
		return wa.readDate(offset)
	}

	def, err := wa.typeInfo.GetTypeDefinition(ctx, typ)
	if err != nil {
		return nil, err
	}

	id, _ := wa.mod.Memory().ReadUint32Le(offset - 8)
	if id != def.Id {
		return nil, fmt.Errorf("pointer is not to a %s", typ)
	}

	if wa.typeInfo.IsListType(typ) {
		return wa.readArray(ctx, typ, offset)
	} else if wa.typeInfo.IsMapType(typ) {
		return wa.readMap(ctx, typ, offset)
	}

	return wa.readClass(ctx, typ, offset)
}

func (wa *wasmAdapter) writeObject(ctx context.Context, typ string, val any) (offset uint32, err error) {
	switch typ {
	case "~lib/arraybuffer/ArrayBuffer":
		switch v := val.(type) {
		case []byte:
			return wa.writeBytes(ctx, v)
		case *[]byte:
			if v == nil {
				return 0, nil
			}
			return wa.writeBytes(ctx, *v)
		default:
			return 0, fmt.Errorf("input value is not a byte array")
		}

	case "string", "~lib/string/String":
		switch v := val.(type) {
		case string:
			return wa.writeString(ctx, v)
		case *string:
			if v == nil {
				return 0, nil
			}
			return wa.writeString(ctx, *v)
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

		return wa.writeDate(ctx, t)
	}

	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		v := reflect.ValueOf(val)
		if v.IsNil() {
			return 0, nil
		} else {
			val = v.Elem().Interface()
		}
	}

	if wa.typeInfo.IsListType(typ) {
		return wa.writeArray(ctx, typ, val)
	} else if wa.typeInfo.IsMapType(typ) {
		return wa.writeMap(ctx, typ, val)
	} else {
		return wa.writeClass(ctx, typ, val)
	}
}
