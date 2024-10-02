/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export function testMapInput_u8_string(m: Map<u8, string>): void {
    assert(m.size == 3, "expected size 3");

    assert(m.has(1), "expected key 1 to be present");
    assert(m.get(1) == "a", "expected value for key 1 to be 'a'");

    assert(m.has(2), "expected key 2 to be present");
    assert(m.get(2) == "b", "expected value for key 2 to be 'b'");

    assert(m.has(3), "expected key 3 to be present");
    assert(m.get(3) == "c", "expected value for key 3 to be 'c'");
}

export function testMapOutput_u8_string(): Map<u8, string> {
    const m = new Map<u8, string>();
    m.set(1, "a");
    m.set(2, "b");
    m.set(3, "c");
    return m;
}

export function testMapInput_string_string(m: Map<string, string>): void {
    assert(m.size == 3, "expected size 3");

    assert(m.has("a"), "expected key 'a' to be present");
    assert(m.get("a") == "1", "expected value for key 'a' to be '1'");

    assert(m.has("b"), "expected key 'b' to be present");
    assert(m.get("b") == "2", "expected value for key 'b' to be '2'");

    assert(m.has("c"), "expected key 'c' to be present");
    assert(m.get("c") == "3", "expected value for key 'c' to be '3'");
}

export function testMapOutput_string_string(): Map<string, string> {
    const m = new Map<string, string>();
    m.set("a", "1");
    m.set("b", "2");
    m.set("c", "3");
    return m;
}

export function testNullableMapInput_string_string(m: Map<string, string> | null): void {
    assert(m != null, "expected non-null map");
    assert(m!.size == 3, "expected size 3");

    assert(m!.has("a"), "expected key 'a' to be present");
    assert(m!.get("a") == "1", "expected value for key 'a' to be '1'");

    assert(m!.has("b"), "expected key 'b' to be present");
    assert(m!.get("b") == "2", "expected value for key 'b' to be '2'");

    assert(m!.has("c"), "expected key 'c' to be present");
    assert(m!.get("c") == "3", "expected value for key 'c' to be '3'");
}

export function testNullableMapInput_string_string_null(m: Map<string, string> | null): void {
    assert(m == null, "expected null map");
}

export function testNullableMapOutput_string_string(): Map<string, string> | null {
    const m = new Map<string, string>();
    m.set("a", "1");
    m.set("b", "2");
    m.set("c", "3");
    return m;
}

export function testNullableMapOutput_string_string_null(): Map<string, string> | null {
    return null;
}

export function testIterateMap_string_string(m: Map<string, string>): void {
    const keys = m.keys();
    const values = m.values();
    for (let i = 0; i < m.size; i++) {
        const key = keys[i];
        const value = values[i];
        console.log(key + " => " + value);
    }
}

export function testMapLookup_string_string(m: Map<string, string>, key: string): string | null {
    if (m.has(key)) {
        return m.get(key);
    }
    return null;
}

class TestClassWithMap {
    m!: Map<string, string>;
}

export function testClassContainingMapInput_string_string(c: TestClassWithMap): void {
    assert(c != null, "expected non-null class");
    assert(c.m != null, "expected non-null map");
    assert(c.m.size == 3, "expected size 3");

    assert(c.m.has("a"), "expected key 'a' to be present");
    assert(c.m.get("a") == "1", "expected value for key 'a' to be '1'");

    assert(c.m.has("b"), "expected key 'b' to be present");
    assert(c.m.get("b") == "2", "expected value for key 'b' to be '2'");

    assert(c.m.has("c"), "expected key 'c' to be present");
    assert(c.m.get("c") == "3", "expected value for key 'c' to be '3'");
}

export function testClassContainingMapOutput_string_string(): TestClassWithMap {
    const c = new TestClassWithMap();
    c.m = new Map<string, string>();
    c.m.set("a", "1");
    c.m.set("b", "2");
    c.m.set("c", "3");
    return c;
}
