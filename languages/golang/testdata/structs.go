package main

type TestStruct1 struct {
	A bool
}

type TestStruct2 struct {
	A bool
	B int
}

type TestStruct3 struct {
	A bool
	B int
	C string
}

type TestStruct4 struct {
	A bool
	B int
	C *string
}

type TestRecursiveStruct struct {
	A bool
	B *TestRecursiveStruct
}

var testStruct1 = TestStruct1{
	A: true,
}

var testStruct2 = TestStruct2{
	A: true,
	B: 123,
}

var testStruct3 = TestStruct3{
	A: true,
	B: 123,
	C: "abc",
}

var testStruct4 = TestStruct4{
	A: true,
	B: 123,
	C: func() *string { s := "abc"; return &s }(),
}

var testStruct4_withNil = TestStruct4{
	A: true,
	B: 123,
	C: nil,
}

var testRecursiveStruct = func() TestRecursiveStruct {
	r := TestRecursiveStruct{
		A: true,
	}
	r.B = &r
	return r
}()

func TestStructInput1(o TestStruct1) {
	assertEqual(testStruct1, o)
}

func TestStructInput2(o TestStruct2) {
	assertEqual(testStruct2, o)
}

func TestStructInput3(o TestStruct3) {
	assertEqual(testStruct3, o)
}

func TestStructInput4(o TestStruct4) {
	assertEqual(testStruct4, o)
}

func TestStructInput4_withNil(o TestStruct4) {
	assertEqual(testStruct4_withNil, o)
}

func TestRecursiveStructInput(o TestRecursiveStruct) {
	assertEqual(testRecursiveStruct, o)
}

func TestStructPtrInput1(o *TestStruct1) {
	assertEqual(testStruct1, *o)
}

func TestStructPtrInput2(o *TestStruct2) {
	assertEqual(testStruct2, *o)
}

func TestStructPtrInput3(o *TestStruct3) {
	assertEqual(testStruct3, *o)
}

func TestStructPtrInput4(o *TestStruct4) {
	assertEqual(testStruct4, *o)
}

func TestStructPtrInput4_withNil(o *TestStruct4) {
	assertEqual(testStruct4_withNil, *o)
}

func TestRecursiveStructPtrInput(o *TestRecursiveStruct) {
	assertEqual(testRecursiveStruct, *o)
}

func TestStructPtrInput1_nil(o *TestStruct1) {
	assertNil(o)
}

func TestStructPtrInput2_nil(o *TestStruct2) {
	assertNil(o)
}

func TestStructPtrInput3_nil(o *TestStruct3) {
	assertNil(o)
}

func TestStructPtrInput4_nil(o *TestStruct4) {
	assertNil(o)
}

func TestRecursiveStructPtrInput_nil(o *TestRecursiveStruct) {
	assertNil(o)
}

func TestStructOutput1() TestStruct1 {
	return testStruct1
}

func TestStructOutput2() TestStruct2 {
	return testStruct2
}

func TestStructOutput3() TestStruct3 {
	return testStruct3
}

func TestStructOutput4() TestStruct4 {
	return testStruct4
}

func TestStructOutput4_withNil() TestStruct4 {
	return testStruct4_withNil
}

func TestRecursiveStructOutput() TestRecursiveStruct {
	return testRecursiveStruct
}

func TestStructPtrOutput1() *TestStruct1 {
	return &testStruct1
}

func TestStructPtrOutput2() *TestStruct2 {
	return &testStruct2
}

func TestStructPtrOutput3() *TestStruct3 {
	return &testStruct3
}

func TestStructPtrOutput4() *TestStruct4 {
	return &testStruct4
}

func TestStructPtrOutput4_withNil() *TestStruct4 {
	return &testStruct4_withNil
}

func TestRecursiveStructPtrOutput() *TestRecursiveStruct {
	return &testRecursiveStruct
}

func TestStructPtrOutput1_nil() *TestStruct1 {
	return nil
}

func TestStructPtrOutput2_nil() *TestStruct2 {
	return nil
}

func TestStructPtrOutput3_nil() *TestStruct3 {
	return nil
}

func TestStructPtrOutput4_nil() *TestStruct4 {
	return nil
}

func TestRecursiveStructPtrOutput_nil() *TestRecursiveStruct {
	return nil
}
