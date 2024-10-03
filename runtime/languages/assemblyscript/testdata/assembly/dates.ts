/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

const testTime = Date.parse("2024-12-31T23:59:59.999Z").getTime();

export function testDateInput(d: Date): void {
    const ts = d.getTime();
    assert(ts == testTime, `Expected ${testTime}, got ${ts}`);
}

export function testDateOutput(): Date {
    return new Date(testTime);
}

export function testNullDateInput(d: Date | null): void {
    assert(d != null, `Expected non-null, got null`);

    const ts = d!.getTime();
    assert(ts == testTime, `Expected ${testTime}, got ${ts}`);
}

export function testNullDateOutput(): Date | null {
    return new Date(testTime);
}

export function testNullDateInput_null(d: Date | null): void {
    assert(d == null, `Expected null, got ${d!.getTime()}`);
}

export function testNullDateOutput_null(): Date | null {
    return null;
}
