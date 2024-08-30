package main

type TestStruct1 struct {
	A bool
}

type TestStruct2 struct {
	A bool
	B int
}

var testStruct1 = TestStruct1{
	A: true,
}

var testStruct2 = TestStruct2{
	A: true,
	B: 123,
}

func TestStructInput1(o TestStruct1) {
	assertEqual(testStruct1, o)
}

func TestStructInput2(o TestStruct2) {
	assertEqual(testStruct2, o)
}

func TestStructPtrInput1(o *TestStruct1) {
	assertEqual(testStruct1, *o)
}

func TestStructPtrInput2(o *TestStruct2) {
	assertEqual(testStruct2, *o)
}

func TestStructPtrInput1_nil(o *TestStruct1) {
	assertNil(o)
}

func TestStructPtrInput2_nil(o *TestStruct2) {
	assertNil(o)
}

func TestStructOutput1() TestStruct1 {
	return TestStruct1{
		A: true,
	}
}

func TestStructOutput2() TestStruct2 {
	return TestStruct2{
		A: true,
		B: 123,
	}
}

func TestStructPtrOutput1() *TestStruct1 {
	return &TestStruct1{
		A: true,
	}
}

func TestStructPtrOutput2() *TestStruct2 {
	return &TestStruct2{
		A: true,
		B: 123,
	}
}

func TestStructPtrOutput1_nil() *TestStruct1 {
	return nil
}

func TestStructPtrOutput2_nil() *TestStruct2 {
	return nil
}
