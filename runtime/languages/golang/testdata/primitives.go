/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import "math"

func TestBoolInput_false(b bool) {
	assertEqual(false, b)
}

func TestBoolInput_true(b bool) {
	assertEqual(true, b)
}

func TestBoolOutput_false() bool {
	return false
}

func TestBoolOutput_true() bool {
	return true
}

func TestBoolPtrInput_false(b *bool) {
	assertEqual(false, *b)
}

func TestBoolPtrInput_true(b *bool) {
	assertEqual(true, *b)
}

func TestBoolPtrInput_nil(b *bool) {
	assertNil(b)
}

func TestBoolPtrOutput_false() *bool {
	b := false
	return &b
}

func TestBoolPtrOutput_true() *bool {
	b := true
	return &b
}

func TestBoolPtrOutput_nil() *bool {
	return nil
}

func TestByteInput_min(b byte) {
	assertEqual(0, b)
}

func TestByteInput_max(b byte) {
	assertEqual(math.MaxUint8, b)
}

func TestByteOutput_min() byte {
	return 0
}

func TestByteOutput_max() byte {
	return math.MaxUint8
}

func TestBytePtrInput_min(b *byte) {
	assertEqual(0, *b)
}

func TestBytePtrInput_max(b *byte) {
	assertEqual(math.MaxUint8, *b)
}

func TestBytePtrInput_nil(b *byte) {
	assertNil(b)
}

func TestBytePtrOutput_min() *byte {
	b := byte(0)
	return &b
}

func TestBytePtrOutput_max() *byte {
	b := byte(math.MaxUint8)
	return &b
}

func TestBytePtrOutput_nil() *byte {
	return nil
}

func TestRuneInput_min(r rune) {
	assertEqual(math.MinInt16, r)
}

func TestRuneInput_max(r rune) {
	assertEqual(math.MaxInt16, r)
}

func TestRuneOutput_min() rune {
	return math.MinInt16
}

func TestRuneOutput_max() rune {
	return math.MaxInt16
}

func TestRunePtrInput_min(r *rune) {
	assertEqual(math.MinInt16, *r)
}

func TestRunePtrInput_max(r *rune) {
	assertEqual(math.MaxInt16, *r)
}

func TestRunePtrInput_nil(r *rune) {
	assertNil(r)
}

func TestRunePtrOutput_min() *rune {
	r := rune(math.MinInt16)
	return &r
}

func TestRunePtrOutput_max() *rune {
	r := rune(math.MaxInt16)
	return &r
}

func TestRunePtrOutput_nil() *rune {
	return nil
}

func TestIntInput_min(n int) {
	assertEqual(math.MinInt32, n)
}

func TestIntInput_max(n int) {
	assertEqual(math.MaxInt32, n)
}

func TestIntOutput_min() int {
	return math.MinInt32
}

func TestIntOutput_max() int {
	return math.MaxInt32
}

func TestIntPtrInput_min(n *int) {
	assertEqual(math.MinInt32, *n)
}

func TestIntPtrInput_max(n *int) {
	assertEqual(math.MaxInt32, *n)
}

func TestIntPtrInput_nil(n *int) {
	assertNil(n)
}

func TestIntPtrOutput_min() *int {
	n := int(math.MinInt32)
	return &n
}

func TestIntPtrOutput_max() *int {
	n := int(math.MaxInt32)
	return &n
}

func TestIntPtrOutput_nil() *int {
	return nil
}

func TestInt8Input_min(n int8) {
	assertEqual(math.MinInt8, n)
}

func TestInt8Input_max(n int8) {
	assertEqual(math.MaxInt8, n)
}

func TestInt8Output_min() int8 {
	return math.MinInt8
}

func TestInt8Output_max() int8 {
	return math.MaxInt8
}

func TestInt8PtrInput_min(n *int8) {
	assertEqual(math.MinInt8, *n)
}

func TestInt8PtrInput_max(n *int8) {
	assertEqual(math.MaxInt8, *n)
}

func TestInt8PtrInput_nil(n *int8) {
	assertNil(n)
}

func TestInt8PtrOutput_min() *int8 {
	n := int8(math.MinInt8)
	return &n
}

func TestInt8PtrOutput_max() *int8 {
	n := int8(math.MaxInt8)
	return &n
}

func TestInt8PtrOutput_nil() *int8 {
	return nil
}

func TestInt16Input_min(n int16) {
	assertEqual(math.MinInt16, n)
}

func TestInt16Input_max(n int16) {
	assertEqual(math.MaxInt16, n)
}

func TestInt16Output_min() int16 {
	return math.MinInt16
}

func TestInt16Output_max() int16 {
	return math.MaxInt16
}

func TestInt16PtrInput_min(n *int16) {
	assertEqual(math.MinInt16, *n)
}

func TestInt16PtrInput_max(n *int16) {
	assertEqual(math.MaxInt16, *n)
}

func TestInt16PtrInput_nil(n *int16) {
	assertNil(n)
}

func TestInt16PtrOutput_min() *int16 {
	n := int16(math.MinInt16)
	return &n
}

func TestInt16PtrOutput_max() *int16 {
	n := int16(math.MaxInt16)
	return &n
}

func TestInt16PtrOutput_nil() *int16 {
	return nil
}

func TestInt32Input_min(n int32) {
	assertEqual(math.MinInt32, n)
}

func TestInt32Input_max(n int32) {
	assertEqual(math.MaxInt32, n)
}

func TestInt32Output_min() int32 {
	return math.MinInt32
}

func TestInt32Output_max() int32 {
	return math.MaxInt32
}

func TestInt32PtrInput_min(n *int32) {
	assertEqual(math.MinInt32, *n)
}

func TestInt32PtrInput_max(n *int32) {
	assertEqual(math.MaxInt32, *n)
}

func TestInt32PtrInput_nil(n *int32) {
	assertNil(n)
}

func TestInt32PtrOutput_min() *int32 {
	n := int32(math.MinInt32)
	return &n
}

func TestInt32PtrOutput_max() *int32 {
	n := int32(math.MaxInt32)
	return &n
}

func TestInt32PtrOutput_nil() *int32 {
	return nil
}

func TestInt64Input_min(n int64) {
	assertEqual(math.MinInt64, n)
}

func TestInt64Input_max(n int64) {
	assertEqual(math.MaxInt64, n)
}

func TestInt64Output_min() int64 {
	return math.MinInt64
}

func TestInt64Output_max() int64 {
	return math.MaxInt64
}

func TestInt64PtrInput_min(n *int64) {
	assertEqual(math.MinInt64, *n)
}

func TestInt64PtrInput_max(n *int64) {
	assertEqual(math.MaxInt64, *n)
}

func TestInt64PtrInput_nil(n *int64) {
	assertNil(n)
}

func TestInt64PtrOutput_min() *int64 {
	n := int64(math.MinInt64)
	return &n
}

func TestInt64PtrOutput_max() *int64 {
	n := int64(math.MaxInt64)
	return &n
}

func TestInt64PtrOutput_nil() *int64 {
	return nil
}

func TestUintInput_min(n uint) {
	assertEqual(0, n)
}

func TestUintInput_max(n uint) {
	assertEqual(math.MaxUint32, n)
}

func TestUintOutput_min() uint {
	return 0
}

func TestUintOutput_max() uint {
	return math.MaxUint32
}

func TestUintPtrInput_min(n *uint) {
	assertEqual(0, *n)
}

func TestUintPtrInput_max(n *uint) {
	assertEqual(math.MaxUint32, *n)
}

func TestUintPtrInput_nil(n *uint) {
	assertNil(n)
}

func TestUintPtrOutput_min() *uint {
	n := uint(0)
	return &n
}

func TestUintPtrOutput_max() *uint {
	n := uint(math.MaxUint32)
	return &n
}

func TestUintPtrOutput_nil() *uint {
	return nil
}

func TestUint8Input_min(n uint8) {
	assertEqual(0, n)
}

func TestUint8Input_max(n uint8) {
	assertEqual(math.MaxUint8, n)
}

func TestUint8Output_min() uint8 {
	return 0
}

func TestUint8Output_max() uint8 {
	return math.MaxUint8
}

func TestUint8PtrInput_min(n *uint8) {
	assertEqual(0, *n)
}

func TestUint8PtrInput_max(n *uint8) {
	assertEqual(math.MaxUint8, *n)
}

func TestUint8PtrInput_nil(n *uint8) {
	assertNil(n)
}

func TestUint8PtrOutput_min() *uint8 {
	n := uint8(0)
	return &n
}

func TestUint8PtrOutput_max() *uint8 {
	n := uint8(math.MaxUint8)
	return &n
}

func TestUint8PtrOutput_nil() *uint8 {
	return nil
}

func TestUint16Input_min(n uint16) {
	assertEqual(0, n)
}

func TestUint16Input_max(n uint16) {
	assertEqual(math.MaxUint16, n)
}

func TestUint16Output_min() uint16 {
	return 0
}

func TestUint16Output_max() uint16 {
	return math.MaxUint16
}

func TestUint16PtrInput_min(n *uint16) {
	assertEqual(0, *n)
}

func TestUint16PtrInput_max(n *uint16) {
	assertEqual(math.MaxUint16, *n)
}

func TestUint16PtrInput_nil(n *uint16) {
	assertNil(n)
}

func TestUint16PtrOutput_min() *uint16 {
	n := uint16(0)
	return &n
}

func TestUint16PtrOutput_max() *uint16 {
	n := uint16(math.MaxUint16)
	return &n
}

func TestUint16PtrOutput_nil() *uint16 {
	return nil
}

func TestUint32Input_min(n uint32) {
	assertEqual(0, n)
}

func TestUint32Input_max(n uint32) {
	assertEqual(math.MaxUint32, n)
}

func TestUint32Output_min() uint32 {
	return 0
}

func TestUint32Output_max() uint32 {
	return math.MaxUint32
}

func TestUint32PtrInput_min(n *uint32) {
	assertEqual(0, *n)
}

func TestUint32PtrInput_max(n *uint32) {
	assertEqual(math.MaxUint32, *n)
}

func TestUint32PtrInput_nil(n *uint32) {
	assertNil(n)
}

func TestUint32PtrOutput_min() *uint32 {
	n := uint32(0)
	return &n
}

func TestUint32PtrOutput_max() *uint32 {
	n := uint32(math.MaxUint32)
	return &n
}

func TestUint32PtrOutput_nil() *uint32 {
	return nil
}

func TestUint64Input_min(n uint64) {
	assertEqual(0, n)
}

func TestUint64Input_max(n uint64) {
	assertEqual(math.MaxUint64, n)
}

func TestUint64Output_min() uint64 {
	return 0
}

func TestUint64Output_max() uint64 {
	return math.MaxUint64
}

func TestUint64PtrInput_min(n *uint64) {
	assertEqual(0, *n)
}

func TestUint64PtrInput_max(n *uint64) {
	assertEqual(math.MaxUint64, *n)
}

func TestUint64PtrInput_nil(n *uint64) {
	assertNil(n)
}

func TestUint64PtrOutput_min() *uint64 {
	n := uint64(0)
	return &n
}

func TestUint64PtrOutput_max() *uint64 {
	n := uint64(math.MaxUint64)
	return &n
}

func TestUint64PtrOutput_nil() *uint64 {
	return nil
}

func TestUintptrInput_min(n uintptr) {
	assertEqual(0, n)
}

func TestUintptrInput_max(n uintptr) {
	assertEqual(math.MaxUint32, n)
}

func TestUintptrOutput_min() uintptr {
	return 0
}

func TestUintptrOutput_max() uintptr {
	return math.MaxUint32
}

func TestUintptrPtrInput_min(n *uintptr) {
	assertEqual(0, *n)
}

func TestUintptrPtrInput_max(n *uintptr) {
	assertEqual(math.MaxUint32, *n)
}

func TestUintptrPtrInput_nil(n *uintptr) {
	assertNil(n)
}

func TestUintptrPtrOutput_min() *uintptr {
	n := uintptr(0)
	return &n
}

func TestUintptrPtrOutput_max() *uintptr {
	n := uintptr(math.MaxUint32)
	return &n
}

func TestUintptrPtrOutput_nil() *uintptr {
	return nil
}

func TestFloat32Input_min(n float32) {
	assertEqual(math.SmallestNonzeroFloat32, n)
}

func TestFloat32Input_max(n float32) {
	assertEqual(math.MaxFloat32, n)
}

func TestFloat32Output_min() float32 {
	return math.SmallestNonzeroFloat32
}

func TestFloat32Output_max() float32 {
	return math.MaxFloat32
}

func TestFloat32PtrInput_min(n *float32) {
	assertEqual(math.SmallestNonzeroFloat32, *n)
}

func TestFloat32PtrInput_max(n *float32) {
	assertEqual(math.MaxFloat32, *n)
}

func TestFloat32PtrInput_nil(n *float32) {
	assertNil(n)
}

func TestFloat32PtrOutput_min() *float32 {
	n := float32(math.SmallestNonzeroFloat32)
	return &n
}

func TestFloat32PtrOutput_max() *float32 {
	n := float32(math.MaxFloat32)
	return &n
}

func TestFloat32PtrOutput_nil() *float32 {
	return nil
}

func TestFloat64Input_min(n float64) {
	assertEqual(math.SmallestNonzeroFloat64, n)
}

func TestFloat64Input_max(n float64) {
	assertEqual(math.MaxFloat64, n)
}

func TestFloat64Output_min() float64 {
	return math.SmallestNonzeroFloat64
}

func TestFloat64Output_max() float64 {
	return math.MaxFloat64
}

func TestFloat64PtrInput_min(n *float64) {
	assertEqual(math.SmallestNonzeroFloat64, *n)
}

func TestFloat64PtrInput_max(n *float64) {
	assertEqual(math.MaxFloat64, *n)
}

func TestFloat64PtrInput_nil(n *float64) {
	assertNil(n)
}

func TestFloat64PtrOutput_min() *float64 {
	n := float64(math.SmallestNonzeroFloat64)
	return &n
}

func TestFloat64PtrOutput_max() *float64 {
	n := float64(math.MaxFloat64)
	return &n
}

func TestFloat64PtrOutput_nil() *float64 {
	return nil
}
