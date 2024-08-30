package main

// "Hello World" in Japanese
const testString = "こんにちは、世界"

func TestStringInput(s string) {
	assertEqual(testString, s)
}

func TestStringPtrInput(s *string) {
	assertEqual(testString, *s)
}

func TestStringPtrInput_nil(s *string) {
	assertNil(s)
}

func TestStringOutput() string {
	return testString
}

func TestStringPtrOutput() *string {
	s := testString
	return &s
}

func TestStringPtrOutput_nil() *string {
	return nil
}
