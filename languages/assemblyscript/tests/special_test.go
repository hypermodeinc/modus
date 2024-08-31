/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript_test

import (
	"fmt"
	"strings"
	"testing"

	"hypruntime/plugins/metadata"

	wasm "github.com/tetratelabs/wazero/api"
)

func TestGetRuntimeTypeInfo(t *testing.T) {

	f := NewASWasmTestFixture(t)
	defer f.Close()

	rtti, err := GetRuntimeTypeInfo(f.Module, f.Plugin.Metadata)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rtti)
}

func GetRuntimeTypeInfo(mod wasm.Module, md *metadata.Metadata) (string, error) {

	types := make(map[uint32]*metadata.TypeDefinition, len(md.Types))
	for _, t := range md.Types {
		types[t.Id] = t
	}

	sb := &strings.Builder{}

	gbl := mod.ExportedGlobal("__rtti_base")
	offset := uint32(gbl.Get())
	sb.WriteString("\nRuntime type information:\n")
	sb.WriteString(fmt.Sprintf("base: 0x%x\n", offset))

	len, ok := mod.Memory().ReadUint32Le(offset)
	offset += 4
	if !ok {
		return "", fmt.Errorf("failed to read length")
	}
	sb.WriteString(fmt.Sprintf("length: %v\n", len))

	sb.WriteString("\n")

	for i := uint32(0); i < len; i++ {
		if i > 0 {
			sb.WriteString("\n")
		}

		var typeName string
		typ := types[i]
		if typ != nil {
			typeName = typ.Name
		} else {
			switch i {
			case 0:
				typeName = "~lib/object/Object"
			case 1:
				typeName = "~lib/arraybuffer/ArrayBuffer"
			case 2:
				typeName = "~lib/string/String"
			default:
				typeName = "(unknown)"
			}
		}

		sb.WriteString(fmt.Sprintf("%v: %s\n", i, typeName))

		f, ok := mod.Memory().ReadUint32Le(offset)
		offset += 4
		if !ok {
			return "", fmt.Errorf("failed to read flags at %v", i)
		}
		flags := RttiFlags(f)

		if flags == RttiNone {
			sb.WriteString("  None\n")
			continue
		}

		if flags&RttiArrayBufferView != 0 {
			sb.WriteString("  ArrayBufferView\n")
		}
		if flags&RttiArray != 0 {
			sb.WriteString("  Array\n")
		}
		if flags&RttiStaticArray != 0 {
			sb.WriteString("  StaticArray\n")
		}
		if flags&RttiSet != 0 {
			sb.WriteString("  Set\n")
		}
		if flags&RttiMap != 0 {
			sb.WriteString("  Map\n")
		}
		if flags&RttiPointerFree != 0 {
			sb.WriteString("  PointerFree\n")
		}
		if flags&RttiValueAlign0 != 0 {
			sb.WriteString("  ValueAlign0\n")
		}
		if flags&RttiValueAlign1 != 0 {
			sb.WriteString("  ValueAlign1\n")
		}
		if flags&RttiValueAlign2 != 0 {
			sb.WriteString("  ValueAlign2\n")
		}
		if flags&RttiValueAlign3 != 0 {
			sb.WriteString("  ValueAlign3\n")
		}
		if flags&RttiValueAlign4 != 0 {
			sb.WriteString("  ValueAlign4\n")
		}
		if flags&RttiValueSigned != 0 {
			sb.WriteString("  ValueSigned\n")
		}
		if flags&RttiValueFloat != 0 {
			sb.WriteString("  ValueFloat\n")
		}
		if flags&RttiValueNullable != 0 {
			sb.WriteString("  ValueNullable\n")
		}
		if flags&RttiValueManaged != 0 {
			sb.WriteString("  ValueManaged\n")
		}
		if flags&RttiKeyAlign0 != 0 {
			sb.WriteString("  KeyAlign0\n")
		}
		if flags&RttiKeyAlign1 != 0 {
			sb.WriteString("  KeyAlign1\n")
		}
		if flags&RttiKeyAlign2 != 0 {
			sb.WriteString("  KeyAlign2\n")
		}
		if flags&RttiKeyAlign3 != 0 {
			sb.WriteString("  KeyAlign3\n")
		}
		if flags&RttiKeyAlign4 != 0 {
			sb.WriteString("  KeyAlign4\n")
		}
		if flags&RttiKeySigned != 0 {
			sb.WriteString("  KeySigned\n")
		}
		if flags&RttiKeyFloat != 0 {
			sb.WriteString("  KeyFloat\n")
		}
		if flags&RttiKeyNullable != 0 {
			sb.WriteString("  KeyNullable\n")
		}
		if flags&RttiKeyManaged != 0 {
			sb.WriteString("  KeyManaged\n")
		}
	}

	return sb.String(), nil
}

// Runtime type information flags.
type RttiFlags uint32

const (
	// No specific flags.
	RttiNone RttiFlags = 0

	// Type is an ArrayBufferView.
	RttiArrayBufferView RttiFlags = 1 << (iota - 1)

	// Type is an Array.
	RttiArray

	// Type is a StaticArray.
	RttiStaticArray

	// Type is a Set.
	RttiSet

	// Type is a Map.
	RttiMap

	// Type has no outgoing pointers.
	RttiPointerFree

	// Value alignment of 1 byte.
	RttiValueAlign0

	// Value alignment of 2 bytes.
	RttiValueAlign1

	// Value alignment of 4 bytes.
	RttiValueAlign2

	// Value alignment of 8 bytes.
	RttiValueAlign3

	// Value alignment of 16 bytes.
	RttiValueAlign4

	// Value is a signed type.
	RttiValueSigned

	// Value is a float type.
	RttiValueFloat

	// Value type is nullable.
	RttiValueNullable

	// Value type is managed.
	RttiValueManaged

	// Key alignment of 1 byte.
	RttiKeyAlign0

	// Key alignment of 2 bytes.
	RttiKeyAlign1

	// Key alignment of 4 bytes.
	RttiKeyAlign2

	// Key alignment of 8 bytes.
	RttiKeyAlign3

	// Key alignment of 16 bytes.
	RttiKeyAlign4

	// Key is a signed type.
	RttiKeySigned

	// Key is a float type.
	RttiKeyFloat

	// Key type is nullable.
	RttiKeyNullable

	// Key type is managed.
	RttiKeyManaged
)
