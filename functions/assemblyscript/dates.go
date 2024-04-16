/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func readDate(mem wasm.Memory, offset uint32) (utils.JSONTime, error) {
	val, ok := mem.ReadUint64Le(offset + 16)
	if !ok {
		return utils.JSONTime{}, fmt.Errorf("error reading timestamp from wasm memory")
	}
	ts := int64(val)
	return utils.JSONTime(time.UnixMilli(ts).UTC()), nil
}

func writeDate(ctx context.Context, mod wasm.Module, t time.Time) (uint32, error) {
	def, err := getTypeDefinition(ctx, "~lib/wasi_date/wasi_Date")
	if err != nil {
		def, err = getTypeDefinition(ctx, "~lib/date/Date")
		if err != nil {
			return 0, err
		}
	}

	t = t.UTC()
	bytes := make([]byte, def.Size)
	binary.LittleEndian.PutUint32(bytes, uint32(t.Year()))
	binary.LittleEndian.PutUint32(bytes[4:], uint32(t.Month()))
	binary.LittleEndian.PutUint32(bytes[8:], uint32(t.Day()))
	binary.LittleEndian.PutUint64(bytes[16:], uint64(t.UnixMilli()))

	return writeObject(ctx, mod, bytes, def.Id), nil
}
