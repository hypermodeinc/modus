/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"errors"

	"hmruntime/functions/assemblyscript"

	wasm "github.com/tetratelabs/wazero/api"
)

func writeResult[T any](ctx context.Context, mod wasm.Module, val T) (uint32, error) {
	switch any(val).(type) {
	case string:
		// fast path for strings
		return assemblyscript.WriteString(ctx, mod, any(val).(string))
	default:
		typ := assemblyscript.GetTypeInfo[T]()
		p, err := assemblyscript.EncodeValue(ctx, mod, typ, val)
		return uint32(p), err
	}
}

func readParam[T any](ctx context.Context, mod wasm.Module, p uint32, v *T) error {
	switch any(*v).(type) {
	case string:
		// fast path for strings
		mem := mod.Memory()
		s, err := assemblyscript.ReadString(mem, p)
		if err != nil {
			return err
		}
		*v = any(s).(T)
	default:
		typ := assemblyscript.GetTypeInfo[T]()
		data, err := assemblyscript.DecodeValueAs[T](ctx, mod, typ, uint64(p))
		if err != nil {
			return err
		}
		*v = any(data).(T)
	}

	return nil
}

func readParams2[T1, T2 any](ctx context.Context, mod wasm.Module,
	p1, p2 uint32,
	v1 *T1, v2 *T2,
) error {
	err1 := readParam[T1](ctx, mod, p1, v1)
	err2 := readParam[T2](ctx, mod, p2, v2)
	return errors.Join(err1, err2)
}

func readParams3[T1, T2, T3 any](ctx context.Context, mod wasm.Module,
	p1, p2, p3 uint32,
	v1 *T1, v2 *T2, v3 *T3,
) error {
	err1 := readParam[T1](ctx, mod, p1, v1)
	err2 := readParam[T2](ctx, mod, p2, v2)
	err3 := readParam[T3](ctx, mod, p3, v3)
	return errors.Join(err1, err2, err3)
}

func readParams4[T1, T2, T3, T4 any](ctx context.Context, mod wasm.Module,
	p1, p2, p3, p4 uint32,
	v1 *T1, v2 *T2, v3 *T3, v4 *T4,
) error {
	err1 := readParam[T1](ctx, mod, p1, v1)
	err2 := readParam[T2](ctx, mod, p2, v2)
	err3 := readParam[T3](ctx, mod, p3, v3)
	err4 := readParam[T4](ctx, mod, p4, v4)
	return errors.Join(err1, err2, err3, err4)
}

func readParams5[T1, T2, T3, T4, T5 any](ctx context.Context, mod wasm.Module,
	p1, p2, p3, p4, p5 uint32,
	v1 *T1, v2 *T2, v3 *T3, v4 *T4, v5 *T5,
) error {
	err1 := readParam[T1](ctx, mod, p1, v1)
	err2 := readParam[T2](ctx, mod, p2, v2)
	err3 := readParam[T3](ctx, mod, p3, v3)
	err4 := readParam[T4](ctx, mod, p4, v4)
	err5 := readParam[T5](ctx, mod, p5, v5)
	return errors.Join(err1, err2, err3, err4, err5)
}

// uncomment to enable as needed (or add more)

// func readParams6[T1, T2, T3, T4, T5, T6 any](ctx context.Context, mod wasm.Module,
// 	p1, p2, p3, p4, p5, p6 uint32,
// 	v1 *T1, v2 *T2, v3 *T3, v4 *T4, v5 *T5, v6 *T6,
// ) error {
// 	err1 := readParam[T1](ctx, mod, p1, v1)
// 	err2 := readParam[T2](ctx, mod, p2, v2)
// 	err3 := readParam[T3](ctx, mod, p3, v3)
// 	err4 := readParam[T4](ctx, mod, p4, v4)
// 	err5 := readParam[T5](ctx, mod, p5, v5)
// 	err6 := readParam[T6](ctx, mod, p6, v6)
// 	return errors.Join(err1, err2, err3, err4, err5, err6)
// }
