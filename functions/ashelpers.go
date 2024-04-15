/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"
	"unicode/utf16"
	"unsafe"

	"hmruntime/plugins"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

// See https://www.assemblyscript.org/runtime.html#memory-layout
type asClassId int64

const (
	asBytes  asClassId = 1
	asString asClassId = 2
)

func writeString(ctx context.Context, mod wasm.Module, s string) uint32 {
	bytes := encodeUTF16(s)
	return writeObject(ctx, mod, bytes, asString)
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

	return writeObject(ctx, mod, bytes, asClassId(def.Id)), nil
}

func writeBytes(ctx context.Context, mod wasm.Module, bytes []byte) uint32 {
	return writeObject(ctx, mod, bytes, asBytes)
}

func writeObject(ctx context.Context, mod wasm.Module, bytes []byte, class asClassId) uint32 {
	offset := allocateWasmMemory(ctx, mod, len(bytes), class)
	mod.Memory().Write(offset, bytes)
	return offset
}

func allocateWasmMemory(ctx context.Context, mod wasm.Module, len int, class asClassId) uint32 {
	// Allocate memory within the AssemblyScript module.
	// This uses the `__new` function exported by the AssemblyScript runtime, so it will be garbage collected.
	// See https://www.assemblyscript.org/runtime.html#interface
	newFn := mod.ExportedFunction("__new")
	res, _ := newFn.Call(ctx, uint64(len), uint64(class))
	return uint32(res[0])
}

func decodeUTF16(bytes []byte) string {

	// Make sure the buffer is valid.
	if len(bytes) == 0 || len(bytes)%2 != 0 {
		return ""
	}

	// Reinterpret []byte as []uint16 to avoid excess copying.
	// This works because we can presume the system is little-endian.
	ptr := unsafe.Pointer(&bytes[0])
	words := unsafe.Slice((*uint16)(ptr), len(bytes)/2)

	// Decode UTF-16 words to a UTF-8 string.
	str := string(utf16.Decode(words))
	return str
}

func encodeUTF16(str string) []byte {
	// Encode the UTF-8 string to UTF-16 words.
	words := utf16.Encode([]rune(str))

	// Reinterpret []uint16 as []byte to avoid excess copying.
	// This works because we can presume the system is little-endian.
	ptr := unsafe.Pointer(&words[0])
	bytes := unsafe.Slice((*byte)(ptr), len(words)*2)
	return bytes
}

func readObject(ctx context.Context, mem wasm.Memory, typ plugins.TypeInfo, offset uint32) (any, error) {
	def, err := getTypeDefinition(ctx, typ.Path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	for _, field := range def.Fields {
		fieldOffset := offset + field.Offset
		result[field.Name], err = readField(ctx, mem, field.Type, fieldOffset)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func readField(ctx context.Context, mem wasm.Memory, typ plugins.TypeInfo, offset uint32) (any, error) {
	var result any
	var ok bool
	switch typ.Name {
	case "bool":
		var val byte
		val, ok = mem.ReadByte(offset)
		result = val != 0

	case "u8":
		result, ok = mem.ReadByte(offset)

	case "u16":
		result, ok = mem.ReadUint16Le(offset)

	case "u32":
		result, ok = mem.ReadUint32Le(offset)

	case "u64":
		result, ok = mem.ReadUint64Le(offset)

	case "i8":
		var val byte
		val, ok = mem.ReadByte(offset)
		result = int8(val)

	case "i16":
		var val uint16
		val, ok = mem.ReadUint16Le(offset)
		result = int16(val)

	case "i32":
		var val uint32
		val, ok = mem.ReadUint32Le(offset)
		result = int32(val)

	case "i64":
		var val uint64
		val, ok = mem.ReadUint64Le(offset)
		result = int64(val)

	case "f32":
		result, ok = mem.ReadFloat32Le(offset)

	case "f64":
		result, ok = mem.ReadFloat64Le(offset)

	default:
		// Managed types have pointers to the actual data.
		p, err := readPointer(mem, offset)
		if err != nil {
			return nil, err
		}

		// Read the actual data.
		switch typ.Name {
		case "string":
			return readString(mem, p)
		case "Date":
			return readDate(mem, p)
		default:
			return readObject(ctx, mem, typ, p)
		}
	}

	if !ok {
		return nil, fmt.Errorf("error reading %s from wasm memory", typ.Name)
	}

	return result, nil
}

func readPointer(mem wasm.Memory, offset uint32) (uint32, error) {
	p, ok := mem.ReadUint32Le(offset)
	if !ok {
		return 0, fmt.Errorf("error reading pointer from wasm memory")
	}

	return p, nil
}

func readDate(mem wasm.Memory, offset uint32) (utils.JSONTime, error) {
	val, ok := mem.ReadUint64Le(offset + 16)
	if !ok {
		return utils.JSONTime{}, fmt.Errorf("error reading timestamp from wasm memory")
	}
	ts := int64(val)
	return utils.JSONTime(time.UnixMilli(ts).UTC()), nil
}

func readString(mem wasm.Memory, offset uint32) (string, error) {

	// AssemblyScript managed objects have their classid stored 8 bytes before the offset.
	// See https://www.assemblyscript.org/runtime.html#memory-layout

	// Read the class id.
	id, ok := mem.ReadUint32Le(offset - 8)
	if !ok {
		return "", fmt.Errorf("failed to read class id of the WASM object")
	}

	// Make sure the pointer is to a string.
	if id != uint32(asString) {
		return "", fmt.Errorf("pointer is not to a string")
	}

	// Read from the buffer and decode it as a string.
	buf, err := readBytes(mem, offset)
	if err != nil {
		return "", err
	}

	return decodeUTF16(buf), nil
}

func readBytes(mem wasm.Memory, offset uint32) ([]byte, error) {

	// The length of AssemblyScript managed objects is stored 4 bytes before the offset.
	// See https://www.assemblyscript.org/runtime.html#memory-layout

	// Read the length.
	len, ok := mem.ReadUint32Le(offset - 4)
	if !ok {
		return nil, fmt.Errorf("failed to read buffer length")
	}

	// Handle empty buffers.
	if len == 0 {
		return []byte{}, nil
	}

	// Now read the data into the buffer.
	buf, ok := mem.Read(offset, len)
	if !ok {
		return nil, fmt.Errorf("failed to read buffer data from WASM memory")
	}

	return buf, nil
}
