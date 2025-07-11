/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

func TestArrayInput0_byte(val [0]byte) {
	expected := [0]byte{}
	assertEqual(expected, val)
}

func TestArrayOutput0_byte() [0]byte {
	return [0]byte{}
}

func TestArrayInput0_string(val [0]string) {
	expected := [0]string{}
	assertEqual(expected, val)
}

func TestArrayOutput0_string() [0]string {
	return [0]string{}
}

func TestArrayInput0_stringPtr(val [0]*string) {
	expected := [0]*string{}
	assertEqual(expected, val)
}

func TestArrayOutput0_stringPtr() [0]*string {
	return [0]*string{}
}

func TestArrayInput0_intPtr(val [0]*int) {
	expected := getIntPtrArray0()
	assertEqual(expected, val)
}

func TestArrayOutput0_intPtr() [0]*int {
	return getIntPtrArray0()
}

func getIntPtrArray0() [0]*int {
	return [0]*int{}
}

func TestArrayInput1_byte(val [1]byte) {
	expected := [1]byte{1}
	assertEqual(expected, val)
}

func TestArrayOutput1_byte() [1]byte {
	return [1]byte{1}
}

func TestArrayInput1_string(val [1]string) {
	expected := [1]string{"abc"}
	assertEqual(expected, val)
}

func TestArrayOutput1_string() [1]string {
	return [1]string{"abc"}
}

func TestArrayInput1_stringPtr(val [1]*string) {
	expected := getStringPtrArray1()
	assertEqual(expected, val)
}

func TestArrayOutput1_stringPtr() [1]*string {
	return getStringPtrArray1()
}

func getStringPtrArray1() [1]*string {
	a := "abc"
	return [1]*string{&a}
}

func TestArrayInput1_intPtr(val [1]*int) {
	expected := getIntPtrArray1()
	assertEqual(expected, val)
}

func TestArrayOutput1_intPtr() [1]*int {
	return getIntPtrArray1()
}

func getIntPtrArray1() [1]*int {
	a := 11
	return [1]*int{&a}
}

func TestArrayInput2_byte(val [2]byte) {
	expected := [2]byte{1, 2}
	assertEqual(expected, val)
}

func TestArrayOutput2_byte() [2]byte {
	return [2]byte{1, 2}
}

func TestArrayInput2_string(val [2]string) {
	expected := [2]string{"abc", "def"}
	assertEqual(expected, val)
}

func TestArrayOutput2_string() [2]string {
	return [2]string{"abc", "def"}
}

func TestArrayInput2_stringPtr(val [2]*string) {
	expected := getStringPtrArray2()
	assertEqual(expected, val)
}

func TestArrayOutput2_stringPtr() [2]*string {
	return getStringPtrArray2()
}

func getStringPtrArray2() [2]*string {
	a := "abc"
	b := "def"
	return [2]*string{&a, &b}
}

func TestArrayInput2_intPtr(val [2]*int) {
	expected := getIntPtrArray2()
	assertEqual(expected, val)
}

func TestArrayOutput2_intPtr() [2]*int {
	return getIntPtrArray2()
}

func getIntPtrArray2() [2]*int {
	a := 11
	b := 22
	return [2]*int{&a, &b}
}

func TestArrayInput2_struct(val [2]TestStruct2) {
	expected := getStructArray2()
	assertEqual(expected, val)
}

func TestArrayOutput2_struct() [2]TestStruct2 {
	return getStructArray2()
}

func TestArrayInput2_structPtr(val [2]*TestStruct2) {
	expected := getStructPtrArray2()
	assertEqual(expected, val)
}

func TestArrayOutput2_structPtr() [2]*TestStruct2 {
	return getStructPtrArray2()
}

func getStructArray2() [2]TestStruct2 {
	return [2]TestStruct2{
		{A: true, B: 123},
		{A: false, B: 456},
	}
}

func getStructPtrArray2() [2]*TestStruct2 {
	return [2]*TestStruct2{
		{A: true, B: 123},
		{A: false, B: 456},
	}
}

func TestArrayInput2_map(val [2]map[string]string) {
	expected := getMapArray2()
	assertEqual(expected, val)
}

func TestArrayOutput2_map() [2]map[string]string {
	return getMapArray2()
}

func TestArrayInput2_mapPtr(val [2]*map[string]string) {
	expected := getMapPtrArray2()
	assertEqual(expected, val)
}

func TestArrayOutput2_mapPtr() [2]*map[string]string {
	return getMapPtrArray2()
}

func getMapArray2() [2]map[string]string {
	return [2]map[string]string{
		{"A": "true", "B": "123"},
		{"C": "false", "D": "456"},
	}
}

func getMapPtrArray2() [2]*map[string]string {
	return [2]*map[string]string{
		{"A": "true", "B": "123"},
		{"C": "false", "D": "456"},
	}
}

func TestPtrArrayInput1_int(val *[1]int) {
	expected := getPtrIntArray1()
	assertEqual(expected, val)
}

func TestPtrArrayOutput1_int() *[1]int {
	return getPtrIntArray1()
}

func getPtrIntArray1() *[1]int {
	a := 11
	return &[1]int{a}
}

func TestPtrArrayInput2_int(val *[2]int) {
	expected := getPtrIntArray2()
	assertEqual(expected, val)
}

func TestPtrArrayOutput2_int() *[2]int {
	return getPtrIntArray2()
}

func getPtrIntArray2() *[2]int {
	a := 11
	b := 22
	return &[2]int{a, b}
}

func TestPtrArrayInput1_string(val *[1]string) {
	expected := getPtrStringArray1()
	assertEqual(expected, val)
}

func TestPtrArrayOutput1_string() *[1]string {
	return getPtrStringArray1()
}

func getPtrStringArray1() *[1]string {
	a := "abc"
	return &[1]string{a}
}

func getPtrStringArray2() *[2]string {
	a := "abc"
	b := "def"
	return &[2]string{a, b}
}

func TestPtrArrayInput2_string(val *[2]string) {
	expected := getPtrStringArray2()
	assertEqual(expected, val)
}

func TestPtrArrayOutput2_string() *[2]string {
	return getPtrStringArray2()
}
