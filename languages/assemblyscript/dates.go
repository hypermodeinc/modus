/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"time"

	"hmruntime/utils"
)

func (wa *wasmAdapter) readDate(offset uint32) (data utils.JSONTime, err error) {
	val, ok := wa.mod.Memory().ReadUint64Le(offset + 16)
	if !ok {
		return utils.JSONTime{}, fmt.Errorf("error reading timestamp from wasm memory")
	}
	ts := int64(val)
	return utils.JSONTime(time.UnixMilli(ts).UTC()), nil
}

func (wa *wasmAdapter) writeDate(ctx context.Context, t time.Time) (offset uint32, err error) {
	typ, err := wa.typeInfo.GetTypeDefinition(ctx, "~lib/wasi_date/wasi_Date")
	if err != nil {
		typ, err = wa.typeInfo.GetTypeDefinition(ctx, "~lib/date/Date")
		if err != nil {
			return 0, err
		}
	}

	const size = 24
	offset, err = wa.allocateWasmMemory(ctx, size, typ.Id)
	if err != nil {
		return 0, err
	}

	t = t.UTC()
	mem := wa.mod.Memory()
	ok1 := mem.WriteUint32Le(offset, uint32(t.Year()))
	ok2 := mem.WriteUint32Le(offset+4, uint32(t.Month()))
	ok3 := mem.WriteUint32Le(offset+8, uint32(t.Day()))
	ok4 := mem.WriteUint64Le(offset+16, uint64(t.UnixMilli()))

	if !(ok1 && ok2 && ok3 && ok4) {
		return 0, fmt.Errorf("failed to write Date object to WASM memory")
	}

	return offset, nil
}
