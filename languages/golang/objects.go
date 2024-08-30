/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"

	"hypruntime/utils"
)

func (wa *wasmAdapter) encodeObject(ctx context.Context, typ string, obj any) ([]uint64, utils.Cleaner, error) {

	if wa.typeInfo.IsPointerType(typ) {
		return wa.encodePointer(ctx, typ, obj)
	}

	if wa.typeInfo.IsStringType(typ) {
		return wa.encodeString(ctx, obj)
	}

	if wa.typeInfo.IsPrimitiveType(typ) {
		return wa.encodePrimitive(typ, obj)
	}

	if wa.typeInfo.IsSliceType(typ) {
		return wa.encodeSlice(ctx, typ, obj)
	}

	if wa.typeInfo.IsArrayType(typ) {
		return wa.encodeArray(ctx, typ, obj)
	}

	if wa.typeInfo.IsMapType(typ) {
		return wa.encodeMap(ctx, typ, obj)
	}

	if wa.typeInfo.IsTimestampType(typ) {
		return wa.encodeTime(ctx, obj)
	}

	return wa.encodeStruct(ctx, typ, obj)
}

func (wa *wasmAdapter) decodeObject(ctx context.Context, typ string, vals []uint64) (any, error) {

	if wa.typeInfo.IsPointerType(typ) {
		return wa.decodePointer(ctx, typ, vals)
	}

	if wa.typeInfo.IsStringType(typ) {
		return wa.decodeString(vals)
	}

	if wa.typeInfo.IsPrimitiveType(typ) {
		return wa.decodePrimitive(typ, vals)
	}

	if wa.typeInfo.IsSliceType(typ) {
		return wa.decodeSlice(ctx, typ, vals)
	}

	if wa.typeInfo.IsArrayType(typ) {
		return wa.decodeArray(ctx, typ, vals)
	}

	if wa.typeInfo.IsMapType(typ) {
		return wa.decodeMap(ctx, typ, vals)
	}

	if wa.typeInfo.IsTimestampType(typ) {
		return wa.decodeTime(vals)
	}

	return wa.decodeStruct(ctx, typ, vals)
}

func (wa *wasmAdapter) writeObject(ctx context.Context, typ string, offset uint32, obj any) (utils.Cleaner, error) {

	if wa.typeInfo.IsPointerType(typ) {
		return wa.writePointer(ctx, typ, offset, obj)
	}

	if wa.typeInfo.IsStringType(typ) {
		return wa.writeString(ctx, offset, obj)
	}

	if wa.typeInfo.IsPrimitiveType(typ) {
		return wa.writePrimitive(typ, offset, obj)
	}

	if wa.typeInfo.IsSliceType(typ) {
		return wa.writeSlice(ctx, typ, offset, obj)
	}

	if wa.typeInfo.IsArrayType(typ) {
		return wa.writeArray(ctx, typ, offset, obj)
	}

	if wa.typeInfo.IsMapType(typ) {
		return wa.writeMap(ctx, typ, offset, obj)
	}

	if wa.typeInfo.IsTimestampType(typ) {
		return wa.writeTime(ctx, offset, obj)
	}

	return wa.writeStruct(ctx, typ, offset, obj)
}

func (wa *wasmAdapter) readObject(ctx context.Context, typ string, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	if wa.typeInfo.IsPointerType(typ) {
		return wa.readPointer(ctx, typ, offset)
	}

	if wa.typeInfo.IsStringType(typ) {
		return wa.readString(offset)
	}

	if wa.typeInfo.IsPrimitiveType(typ) {
		return wa.readPrimitive(typ, offset)
	}

	if wa.typeInfo.IsSliceType(typ) {
		return wa.readSlice(ctx, typ, offset)
	}

	if wa.typeInfo.IsArrayType(typ) {
		return wa.readArray(ctx, typ, offset)
	}

	if wa.typeInfo.IsMapType(typ) {
		return wa.readMap(ctx, typ, offset)
	}

	if wa.typeInfo.IsTimestampType(typ) {
		return wa.readTime(offset)
	}

	return wa.readStruct(ctx, typ, offset)
}
