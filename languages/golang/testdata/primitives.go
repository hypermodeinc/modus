package main

func TestBoolInput(b bool) {
	assertEqual(true, b)
}

func TestBoolPtrInput(b *bool) {
	assertEqual(true, *b)
}

func TestBoolPtrInput_nil(b *bool) {
	assertNil(b)
}

func TestBoolOutput() bool {
	return true
}

func TestBoolPtrOutput() *bool {
	b := true
	return &b
}

func TestBoolPtrOutput_nil() *bool {
	return nil
}

func TestByteInput(b byte) {
	assertEqual(123, b)
}

func TestBytePtrInput(b *byte) {
	assertEqual(123, *b)
}

func TestBytePtrInput_nil(b *byte) {
	assertNil(b)
}

func TestByteOutput() byte {
	return 123
}

func TestBytePtrOutput() *byte {
	b := byte(123)
	return &b
}

func TestBytePtrOutput_nil() *byte {
	return nil
}

func TestFloat32Input(n float32) {
	assertEqual(123, n)
}

func TestFloat32PtrInput(n *float32) {
	assertEqual(123, *n)
}

func TestFloat32PtrInput_nil(n *float32) {
	assertNil(n)
}

func TestFloat32Output() float32 {
	return 123
}

func TestFloat32PtrOutput() *float32 {
	n := float32(123)
	return &n
}

func TestFloat32PtrOutput_nil() *float32 {
	return nil
}

func TestFloat64Input(n float64) {
	assertEqual(123, n)
}

func TestFloat64PtrInput(n *float64) {
	assertEqual(123, *n)
}

func TestFloat64PtrInput_nil(n *float64) {
	assertNil(n)
}

func TestFloat64Output() float64 {
	return 123
}

func TestFloat64PtrOutput() *float64 {
	n := float64(123)
	return &n
}

func TestFloat64PtrOutput_nil() *float64 {
	return nil
}

func TestIntInput(n int) {
	assertEqual(123, n)
}

func TestIntPtrInput(n *int) {
	assertEqual(123, *n)
}

func TestIntPtrInput_nil(n *int) {
	assertNil(n)
}

func TestIntOutput() int {
	return 123
}

func TestIntPtrOutput() *int {
	n := int(123)
	return &n
}

func TestIntPtrOutput_nil() *int {
	return nil
}

func TestInt8Input(n int8) {
	assertEqual(123, n)
}

func TestInt8PtrInput(n *int8) {
	assertEqual(123, *n)
}

func TestInt8PtrInput_nil(n *int8) {
	assertNil(n)
}

func TestInt8Output() int8 {
	return 123
}

func TestInt8PtrOutput() *int8 {
	n := int8(123)
	return &n
}

func TestInt8PtrOutput_nil() *int8 {
	return nil
}

func TestInt16Input(n int16) {
	assertEqual(123, n)
}

func TestInt16PtrInput(n *int16) {
	assertEqual(123, *n)
}

func TestInt16PtrInput_nil(n *int16) {
	assertNil(n)
}

func TestInt16Output() int16 {
	return 123
}

func TestInt16PtrOutput() *int16 {
	n := int16(123)
	return &n
}

func TestInt16PtrOutput_nil() *int16 {
	return nil
}

func TestInt32Input(n int32) {
	assertEqual(123, n)
}

func TestInt32PtrInput(n *int32) {
	assertEqual(123, *n)
}

func TestInt32PtrInput_nil(n *int32) {
	assertNil(n)
}

func TestInt32Output() int32 {
	return 123
}

func TestInt32PtrOutput() *int32 {
	n := int32(123)
	return &n
}

func TestInt32PtrOutput_nil() *int32 {
	return nil
}

func TestInt64Input(n int64) {
	assertEqual(123, n)
}

func TestInt64PtrInput(n *int64) {
	assertEqual(123, *n)
}

func TestInt64PtrInput_nil(n *int64) {
	assertNil(n)
}

func TestInt64Output() int64 {
	return 123
}

func TestInt64PtrOutput() *int64 {
	n := int64(123)
	return &n
}

func TestInt64PtrOutput_nil() *int64 {
	return nil
}

func TestRuneInput(r rune) {
	assertEqual(123, r)
}

func TestRunePtrInput(r *rune) {
	assertEqual(123, *r)
}

func TestRunePtrInput_nil(r *rune) {
	assertNil(r)
}

func TestRuneOutput() rune {
	return 123
}

func TestRunePtrOutput() *rune {
	r := rune(123)
	return &r
}

func TestRunePtrOutput_nil() *rune {
	return nil
}

func TestUintInput(n uint) {
	assertEqual(123, n)
}

func TestUintPtrInput(n *uint) {
	assertEqual(123, *n)
}

func TestUintPtrInput_nil(n *uint) {
	assertNil(n)
}

func TestUintOutput() uint {
	return 123
}

func TestUintPtrOutput() *uint {
	n := uint(123)
	return &n
}

func TestUintPtrOutput_nil() *uint {
	return nil
}

func TestUint8Input(n uint8) {
	assertEqual(123, n)
}

func TestUint8PtrInput(n *uint8) {
	assertEqual(123, *n)
}

func TestUint8PtrInput_nil(n *uint8) {
	assertNil(n)
}

func TestUint8Output() uint8 {
	return 123
}

func TestUint8PtrOutput() *uint8 {
	n := uint8(123)
	return &n
}

func TestUint8PtrOutput_nil() *uint8 {
	return nil
}

func TestUint16Input(n uint16) {
	assertEqual(123, n)
}

func TestUint16PtrInput(n *uint16) {
	assertEqual(123, *n)
}

func TestUint16PtrInput_nil(n *uint16) {
	assertNil(n)
}

func TestUint16Output() uint16 {
	return 123
}

func TestUint16PtrOutput() *uint16 {
	n := uint16(123)
	return &n
}

func TestUint16PtrOutput_nil() *uint16 {
	return nil
}

func TestUint32Input(n uint32) {
	assertEqual(123, n)
}

func TestUint32PtrInput(n *uint32) {
	assertEqual(123, *n)
}

func TestUint32PtrInput_nil(n *uint32) {
	assertNil(n)
}

func TestUint32Output() uint32 {
	return 123
}

func TestUint32PtrOutput() *uint32 {
	n := uint32(123)
	return &n
}

func TestUint32PtrOutput_nil() *uint32 {
	return nil
}

func TestUint64Input(n uint64) {
	assertEqual(123, n)
}

func TestUint64PtrInput(n *uint64) {
	assertEqual(123, *n)
}

func TestUint64PtrInput_nil(n *uint64) {
	assertNil(n)
}

func TestUint64Output() uint64 {
	return 123
}

func TestUint64PtrOutput() *uint64 {
	n := uint64(123)
	return &n
}

func TestUint64PtrOutput_nil() *uint64 {
	return nil
}

func TestUintptrInput(n uintptr) {
	assertEqual(123, n)
}

func TestUintptrPtrInput(n *uintptr) {
	assertEqual(123, *n)
}

func TestUintptrPtrInput_nil(n *uintptr) {
	assertNil(n)
}

func TestUintptrOutput() uintptr {
	return 123
}

func TestUintptrPtrOutput() *uintptr {
	n := uintptr(123)
	return &n
}

func TestUintptrPtrOutput_nil() *uintptr {
	return nil
}
