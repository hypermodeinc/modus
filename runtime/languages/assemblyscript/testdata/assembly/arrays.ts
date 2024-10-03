/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export function testArrayOutput_i8(): i8[] {
    return [1, 2, 3];
}

export function testArrayInput_i8(arr: i8[]): void {
    assert(arr.length == 3);
    assert(arr[0] == 1);
    assert(arr[1] == 2);
    assert(arr[2] == 3);
}

export function testArrayOutput_i8_empty(): i8[] {
    return [];
}

export function testArrayInput_i8_empty(arr: i8[]): void {
    assert(arr.length == 0);
}

export function testArrayOutput_i8_null(): i8[] | null {
    return null;
}

export function testArrayInput_i8_null(arr: i8[] | null): void {
    assert(arr == null);
}

export function testArrayInput_i32(arr: i32[]): void {
    assert(arr.length == 3);
    assert(arr[0] == 1);
    assert(arr[1] == 2);
    assert(arr[2] == 3);
}

export function testArrayOutput_i32(): i32[] {
    return [1, 2, 3];
}

export function testArrayInput_f32(arr: f32[]): void {
    assert(arr.length == 3);
    assert(arr[0] == 1);
    assert(arr[1] == 2);
    assert(arr[2] == 3);
}

export function testArrayOutput_f32(): f32[] {
    return [1, 2, 3];
}

export function testArrayInput_string(arr: string[]): void {
    assert(arr.length == 3);
    assert(arr[0] == "abc");
    assert(arr[1] == "def");
    assert(arr[2] == "ghi");
}

export function testArrayOutput_string(): string[] {
    return ["abc", "def", "ghi"];
}

export function testArrayInput_string_2d(arr: string[][]): void {
    assert(arr.length == 3);
    assert(arr[0].length == 3);
    assert(arr[0][0] == "abc");
    assert(arr[0][1] == "def");
    assert(arr[0][2] == "ghi");
    assert(arr[1].length == 3);
    assert(arr[1][0] == "jkl");
    assert(arr[1][1] == "mno");
    assert(arr[1][2] == "pqr");
    assert(arr[2].length == 3);
    assert(arr[2][0] == "stu");
    assert(arr[2][1] == "vwx");
    assert(arr[2][2] == "yz");
}

export function testArrayOutput_string_2d(): string[][] {
    return [["abc", "def", "ghi"], ["jkl", "mno", "pqr"], ["stu", "vwx", "yz"]];
}

export function testArrayInput_string_2d_empty(arr: string[][]): void {
    assert(arr.length == 0);
}

export function testArrayOutput_string_2d_empty(): string[][] {
    return [];
}

class TestObject1 {
    constructor(public a: i32, public b: i32) { }
}

export function testArrayIteration(arr: TestObject1[]): void {
    for (let i = 0; i < arr.length; i++) {
        let obj = arr[i];
        console.log(`[${i}]: a=${obj.a}, b=${obj.b}`);
    }
}
