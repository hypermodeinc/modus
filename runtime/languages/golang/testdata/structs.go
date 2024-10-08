/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

type TestStruct1 struct {
	A bool
}

type TestStruct1_map struct {
	A bool
}

type TestStruct2 struct {
	A bool
	B int
}

type TestStruct2_map struct {
	A bool
	B int
}

type TestStruct3 struct {
	A bool
	B int
	C string
}

type TestStruct3_map struct {
	A bool
	B int
	C string
}

type TestStruct4 struct {
	A bool
	B int
	C *string
}

type TestStruct4_map struct {
	A bool
	B int
	C *string
}

type TestStruct5 struct {
	A string
	B string
	C string
	D []string
	E float64
	F float64
}

type TestStruct5_map struct {
	A string
	B string
	C string
	D []string
	E float64
	F float64
}

type TestRecursiveStruct struct {
	A bool
	B *TestRecursiveStruct
}

type TestRecursiveStruct_map struct {
	A bool
	B *TestRecursiveStruct_map
}

var testStruct1 = TestStruct1{
	A: true,
}

var testStruct1_map = TestStruct1_map{
	A: true,
}

var testStruct2 = TestStruct2{
	A: true,
	B: 123,
}

var testStruct2_map = TestStruct2_map{
	A: true,
	B: 123,
}

var testStruct3 = TestStruct3{
	A: true,
	B: 123,
	C: "abc",
}

var testStruct3_map = TestStruct3_map{
	A: true,
	B: 123,
	C: "abc",
}

var testStruct4 = TestStruct4{
	A: true,
	B: 123,
	C: func() *string { s := "abc"; return &s }(),
}

var testStruct4_map = TestStruct4_map{
	A: true,
	B: 123,
	C: func() *string { s := "abc"; return &s }(),
}

var testStruct4_withNil = TestStruct4{
	A: true,
	B: 123,
	C: nil,
}

var testStruct4_map_withNil = TestStruct4_map{
	A: true,
	B: 123,
	C: nil,
}

var testStruct5 = TestStruct5{
	A: "abc",
	B: "def",
	C: "ghi",
	D: []string{
		"jkl",
		"mno",
		"pqr",
	},
	E: 0.12345,
	F: 99.99999,
}

var testStruct5_map = TestStruct5_map{
	A: "abc",
	B: "def",
	C: "ghi",
	D: []string{
		"jkl",
		"mno",
		"pqr",
	},
	E: 0.12345,
	F: 99.99999,
}

var testRecursiveStruct = func() TestRecursiveStruct {
	r := TestRecursiveStruct{
		A: true,
	}
	r.B = &r
	return r
}()

var testRecursiveStruct_map = func() TestRecursiveStruct_map {
	r := TestRecursiveStruct_map{
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

func TestStructInput5(o TestStruct5) {
	assertEqual(testStruct5, o)
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

func TestStructPtrInput5(o *TestStruct5) {
	assertEqual(testStruct5, *o)
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

func TestStructPtrInput5_nil(o *TestStruct5) {
	assertNil(o)
}

func TestRecursiveStructPtrInput_nil(o *TestRecursiveStruct) {
	assertNil(o)
}

func TestStructOutput1() TestStruct1 {
	return testStruct1
}

func TestStructOutput1_map() TestStruct1_map {
	return testStruct1_map
}

func TestStructOutput2() TestStruct2 {
	return testStruct2
}

func TestStructOutput2_map() TestStruct2_map {
	return testStruct2_map
}

func TestStructOutput3() TestStruct3 {
	return testStruct3
}

func TestStructOutput3_map() TestStruct3_map {
	return testStruct3_map
}

func TestStructOutput4() TestStruct4 {
	return testStruct4
}

func TestStructOutput4_map() TestStruct4_map {
	return testStruct4_map
}

func TestStructOutput5() TestStruct5 {
	return testStruct5
}

func TestStructOutput5_map() TestStruct5_map {
	return testStruct5_map
}

func TestStructOutput4_withNil() TestStruct4 {
	return testStruct4_withNil
}

func TestStructOutput4_map_withNil() TestStruct4_map {
	return testStruct4_map_withNil
}

func TestRecursiveStructOutput() TestRecursiveStruct {
	return testRecursiveStruct
}

func TestRecursiveStructOutput_map() TestRecursiveStruct_map {
	return testRecursiveStruct_map
}

func TestStructPtrOutput1() *TestStruct1 {
	return &testStruct1
}

func TestStructPtrOutput1_map() *TestStruct1_map {
	return &testStruct1_map
}

func TestStructPtrOutput2() *TestStruct2 {
	return &testStruct2
}

func TestStructPtrOutput2_map() *TestStruct2_map {
	return &testStruct2_map
}

func TestStructPtrOutput3() *TestStruct3 {
	return &testStruct3
}

func TestStructPtrOutput3_map() *TestStruct3_map {
	return &testStruct3_map
}

func TestStructPtrOutput4() *TestStruct4 {
	return &testStruct4
}

func TestStructPtrOutput4_map() *TestStruct4_map {
	return &testStruct4_map
}

func TestStructPtrOutput5() *TestStruct5 {
	return &testStruct5
}

func TestStructPtrOutput5_map() *TestStruct5_map {
	return &testStruct5_map
}

func TestStructPtrOutput4_withNil() *TestStruct4 {
	return &testStruct4_withNil
}

func TestStructPtrOutput4_map_withNil() *TestStruct4_map {
	return &testStruct4_map_withNil
}

func TestRecursiveStructPtrOutput() *TestRecursiveStruct {
	return &testRecursiveStruct
}

func TestRecursiveStructPtrOutput_map() *TestRecursiveStruct_map {
	return &testRecursiveStruct_map
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

func TestStructPtrOutput5_nil() *TestStruct5 {
	return nil
}

func TestRecursiveStructPtrOutput_nil() *TestRecursiveStruct {
	return nil
}
