package main

func TestSliceInput_byte(val []byte) {
	var expected = []byte{1, 2, 3, 4}
	assertSlicesEqual(expected, val)
}

func TestSliceOutput_byte() []byte {
	return []byte{1, 2, 3, 4}
}

func TestSliceInput_string(val []string) {
	var expected = []string{"abc", "def", "ghi"}
	assertSlicesEqual(expected, val)
}

func TestSliceOutput_string() []string {
	return []string{"abc", "def", "ghi"}
}

func TestSliceInput_string_empty(val []string) {
	var expected = []string{}
	assertSlicesEqual(expected, val)
}

func TestSliceOutput_string_empty() []string {
	return []string{}
}

func TestSliceInput_int32_empty(val []int32) {
	var expected = []int32{}
	assertSlicesEqual(expected, val)
}

func TestSliceOutput_int32_empty() []int32 {
	return []int32{}
}

func TestSliceInput_stringPtr(val []*string) {
	var expected = getStringPtrSlice()
	assertPtrSlicesEqual(expected, val)
}

func TestSliceOutput_stringPtr() []*string {
	return getStringPtrSlice()
}

func getStringPtrSlice() []*string {
	a := "abc"
	b := "def"
	c := "ghi"
	return []*string{&a, &b, &c}
}

func TestSliceInput_intPtr(val []*int) {
	var expected = getIntPtrSlice()
	assertPtrSlicesEqual(expected, val)
}

func TestSliceOutput_intPtr() []*int {
	return getIntPtrSlice()
}

func getIntPtrSlice() []*int {
	a := 11
	b := 22
	c := 33
	return []*int{&a, &b, &c}
}
