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

func writeParam[T any](ctx context.Context, mod wasm.Module, val T) (uint32, error) {
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

func readParam[T any](ctx context.Context, mod wasm.Module, p uint32) (T, error) {
	var v T
	switch any(v).(type) {
	case string:
		// fast path for strings
		mem := mod.Memory()
		s, err := assemblyscript.ReadString(mem, p)
		return any(s).(T), err
	default:
		typ := assemblyscript.GetTypeInfo[T]()
		return assemblyscript.DecodeValueAs[T](ctx, mod, typ, uint64(p))
	}
}

func readParams2[T1, T2 any](ctx context.Context, mod wasm.Module, p1, p2 uint32) (T1, T2, error) {
	r1, err1 := readParam[T1](ctx, mod, p1)
	r2, err2 := readParam[T2](ctx, mod, p2)
	return r1, r2, errors.Join(err1, err2)
}

func readParams3[T1, T2, T3 any](ctx context.Context, mod wasm.Module, p1, p2, p3 uint32) (T1, T2, T3, error) {
	r1, err1 := readParam[T1](ctx, mod, p1)
	r2, err2 := readParam[T2](ctx, mod, p2)
	r3, err3 := readParam[T3](ctx, mod, p3)
	return r1, r2, r3, errors.Join(err1, err2, err3)
}

func readParams4[T1, T2, T3, T4 any](ctx context.Context, mod wasm.Module, p1, p2, p3, p4 uint32) (T1, T2, T3, T4, error) {
	r1, err1 := readParam[T1](ctx, mod, p1)
	r2, err2 := readParam[T2](ctx, mod, p2)
	r3, err3 := readParam[T3](ctx, mod, p3)
	r4, err4 := readParam[T4](ctx, mod, p4)
	return r1, r2, r3, r4, errors.Join(err1, err2, err3, err4)
}

// uncomment to enable as needed (or add more)

// func readParams5[T1, T2, T3, T4, T5 any](ctx context.Context, mod wasm.Module, p1, p2, p3, p4, p5 uint32) (T1, T2, T3, T4, T5, error) {
// 	r1, err1 := readParam[T1](ctx, mod, p1)
// 	r2, err2 := readParam[T2](ctx, mod, p2)
// 	r3, err3 := readParam[T3](ctx, mod, p3)
// 	r4, err4 := readParam[T4](ctx, mod, p4)
// 	r5, err5 := readParam[T5](ctx, mod, p5)
// 	return r1, r2, r3, r4, r5, errors.Join(err1, err2, err3, err4, err5)
// }

// func readParam6[T1, T2, T3, T4, T5, T6 any](ctx context.Context, mod wasm.Module, p1, p2, p3, p4, p5, p6 uint32) (T1, T2, T3, T4, T5, T6, error) {
// 	r1, err1 := readParam[T1](ctx, mod, p1)
// 	r2, err2 := readParam[T2](ctx, mod, p2)
// 	r3, err3 := readParam[T3](ctx, mod, p3)
// 	r4, err4 := readParam[T4](ctx, mod, p4)
// 	r5, err5 := readParam[T5](ctx, mod, p5)
// 	r6, err6 := readParam[T6](ctx, mod, p6)
// 	return r1, r2, r3, r4, r5, r6, errors.Join(err1, err2, err3, err4, err5, err6)
// }
