/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"fmt"
	"hypruntime/utils"
	"math"
	"time"

	"github.com/spf13/cast"
)

func (wa *wasmAdapter) encodePrimitive(typ string, obj any) ([]uint64, utils.Cleaner, error) {
	val, err := wa.doEncodePrimitive(typ, obj)
	if err != nil {
		return nil, nil, err
	}
	return []uint64{val}, nil, nil
}

func (wa *wasmAdapter) decodePrimitive(typ string, vals []uint64) (any, error) {
	val := vals[0]

	switch typ {
	case "bool":
		return val != 0, nil

	case "byte":
		return byte(val), nil

	case "rune":
		return rune(val), nil

	case "uint":
		return uint(uint32(val)), nil // uint is 32-bit on wasm

	case "uint8":
		return uint8(val), nil

	case "uint16":
		return uint16(val), nil

	case "uint32":
		return uint32(val), nil

	case "uint64":
		return val, nil

	case "uintptr":
		return uintptr(uint32(val)), nil // uintptr is 32-bit on wasm

	case "int":
		return int(int32(val)), nil // int is 32-bit on wasm

	case "int8":
		return int8(val), nil

	case "int16":
		return int16(val), nil

	case "int32":
		return int32(val), nil

	case "int64":
		return int64(val), nil

	case "time.Duration":
		return time.Duration(val), nil

	case "float32":
		return math.Float32frombits(uint32(val)), nil

	case "float64":
		return math.Float64frombits(val), nil
	}

	return nil, fmt.Errorf("unsupported primitive type: %s", typ)
}

func (wa *wasmAdapter) writePrimitive(typ string, offset uint32, obj any) (utils.Cleaner, error) {
	err := wa.doWritePrimitive(typ, offset, obj)
	return nil, err
}

func (wa *wasmAdapter) readPrimitive(typ string, offset uint32) (any, error) {

	mem := wa.mod.Memory()

	switch typ {
	case "bool":
		val, ok := mem.ReadByte(offset)
		if !ok {
			return false, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return val != 0, nil

	case "byte":
		val, ok := mem.ReadByte(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return val, nil

	case "uint8":
		val, ok := mem.ReadByte(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return uint8(val), nil

	case "uint16":
		val, ok := mem.ReadUint16Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return val, nil

	case "uint32":
		val, ok := mem.ReadUint32Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return val, nil

	case "uint":
		val, ok := mem.ReadUint32Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return uint(val), nil

	case "uintptr":
		val, ok := mem.ReadUint32Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return uintptr(val), nil

	case "uint64":
		val, ok := mem.ReadUint64Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return val, nil

	case "int8":
		val, ok := mem.ReadByte(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return int8(val), nil

	case "int16":
		val, ok := mem.ReadUint16Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return int16(val), nil

	case "int32":
		val, ok := mem.ReadUint32Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return int32(val), nil

	case "rune":
		val, ok := mem.ReadUint32Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return rune(val), nil

	case "int":
		val, ok := mem.ReadUint32Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return int(val), nil

	case "int64":
		val, ok := mem.ReadUint64Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return int64(val), nil

	case "time.Duration":
		val, ok := mem.ReadUint64Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return time.Duration(val), nil

	case "float32":
		val, ok := mem.ReadFloat32Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return val, nil

	case "float64":
		val, ok := mem.ReadFloat64Le(offset)
		if !ok {
			return 0, fmt.Errorf("failed to read %s value from memory", typ)
		}
		return val, nil
	}

	return nil, fmt.Errorf("unsupported primitive type: %s", typ)
}

func (wa *wasmAdapter) doEncodePrimitive(typ string, val any) (uint64, error) {
	switch typ {
	case "bool":
		val, err := cast.ToBoolE(val)
		if err != nil {
			return 0, err
		}

		if val {
			return 1, nil
		} else {
			return 0, nil
		}

	case "uint8", "byte":
		val, err := cast.ToUint8E(val)
		if err != nil {
			return 0, err
		}
		return uint64(val), nil

	case "uint16":
		val, err := cast.ToUint16E(val)
		if err != nil {
			return 0, err
		}
		return uint64(val), nil

	case "uint32", "uint":
		val, err := cast.ToUint32E(val)
		if err != nil {
			return 0, err
		}
		return uint64(val), nil

	case "uint64":
		val, err := cast.ToUint64E(val)
		if err != nil {
			return 0, err
		}
		return val, nil

	case "int8":
		val, err := cast.ToInt8E(val)
		if err != nil {
			return 0, err
		}
		return uint64(uint32(val)), nil

	case "int16":
		val, err := cast.ToInt16E(val)
		if err != nil {
			return 0, err
		}
		return uint64(uint32(val)), nil

	case "int32", "int", "rune":
		val, err := cast.ToInt32E(val)
		if err != nil {
			return 0, err
		}
		return uint64(uint32(val)), nil

	case "int64":
		val, err := cast.ToInt64E(val)
		if err != nil {
			return 0, err
		}
		return uint64(val), nil

	case "time.Duration":
		if v, ok := val.(time.Duration); ok {
			return uint64(v), nil
		}
		return wa.doEncodePrimitive("int64", val)

	case "float32":
		val, err := cast.ToFloat32E(val)
		if err != nil {
			return 0, err
		}
		return uint64(math.Float32bits(val)), nil

	case "float64":
		val, err := cast.ToFloat64E(val)
		if err != nil {
			return 0, err
		}
		return math.Float64bits(val), nil

	case "uintptr":
		if val, ok := val.(uintptr); ok {
			return uint64(uint32(val)), nil
		}
		val, err := cast.ToUint32E(val)
		if err != nil {
			return 0, err
		}
		return uint64(val), nil
	}

	return 0, fmt.Errorf("unsupported primitive type: %s", typ)
}

func (wa *wasmAdapter) doWritePrimitive(typ string, offset uint32, val any) error {
	mem := wa.mod.Memory()

	switch typ {

	case "bool":
		val, err := cast.ToBoolE(val)
		if err != nil {
			return err
		}
		var b byte
		if val {
			b = 1
		}
		if ok := mem.WriteByte(offset, b); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "byte", "uint8":
		val, err := cast.ToUint8E(val)
		if err != nil {
			return err
		}
		if ok := mem.WriteByte(offset, val); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "uint16":
		val, err := cast.ToUint16E(val)
		if err != nil {
			return err
		}
		if ok := mem.WriteUint16Le(offset, val); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "uint", "uint32":
		val, err := cast.ToUint32E(val)
		if err != nil {
			return err
		}
		if ok := mem.WriteUint32Le(offset, val); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "uint64":
		val, err := cast.ToUint64E(val)
		if err != nil {
			return err
		}
		if ok := mem.WriteUint64Le(offset, val); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "int8":
		val, err := cast.ToInt8E(val)
		if err != nil {
			return err
		}
		if ok := mem.WriteByte(offset, byte(val)); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "int16":
		val, err := cast.ToInt16E(val)
		if err != nil {
			return err
		}
		if ok := mem.WriteUint16Le(offset, uint16(val)); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "int", "int32", "rune":
		val, err := cast.ToInt32E(val)
		if err != nil {
			return err
		}
		if ok := mem.WriteUint32Le(offset, uint32(val)); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "int64":
		val, err := cast.ToInt64E(val)
		if err != nil {
			return err
		}
		if ok := mem.WriteUint64Le(offset, uint64(val)); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "time.Duration":
		if v, ok := val.(time.Duration); ok {
			return wa.doWritePrimitive("int64", offset, int64(v))
		}
		return wa.doWritePrimitive("int64", offset, val)

	case "float32":
		val, err := cast.ToFloat32E(val)
		if err != nil {
			return err
		}
		if ok := mem.WriteFloat32Le(offset, val); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "float64":
		val, err := cast.ToFloat64E(val)
		if err != nil {
			return err
		}
		if ok := mem.WriteFloat64Le(offset, val); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	case "uintptr":
		var x uint32
		if v, ok := val.(uintptr); ok {
			x = uint32(v)
		} else if v, err := cast.ToUint32E(val); err == nil {
			x = v
		} else {
			return err
		}

		if ok := mem.WriteUint32Le(offset, x); !ok {
			return fmt.Errorf("failed to write %s value to memory", typ)
		}

	default:
		return fmt.Errorf("unsupported primitive type: %s", typ)
	}

	return nil
}
