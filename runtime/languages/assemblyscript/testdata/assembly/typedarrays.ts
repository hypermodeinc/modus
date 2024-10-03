/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export function testInt8ArrayInput(view: Int8Array): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testInt8ArrayOutput(): Int8Array {
    const view = new Int8Array(4);
    view.set([0, 1, 2, 3]);
    return view;
}

export function testInt16ArrayInput(view: Int16Array): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testInt16ArrayOutput(): Int16Array {
    const view = new Int16Array(4);
    view.set([0, 1, 2, 3]);
    return view;
}

export function testInt32ArrayInput(view: Int32Array): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testInt32ArrayOutput(): Int32Array {
    const view = new Int32Array(4);
    view.set([0, 1, 2, 3]);
    return view;
}

export function testInt64ArrayInput(view: Int64Array): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testInt64ArrayOutput(): Int64Array {
    const view = new Int64Array(4);
    view.set([0, 1, 2, 3]);
    return view;
}

export function testUint8ArrayInput(view: Uint8Array): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testUint8ArrayInput_empty(view: Uint8Array): void {
    assert(view.length == 0);
}

export function testUint8ArrayInput_null(view: Uint8Array | null): void {
    assert(view == null);
}

export function testUint8ArrayOutput(): Uint8Array {
    const view = new Uint8Array(4);
    view.set([0, 1, 2, 3]);
    return view;
}

export function testUint8ArrayOutput_empty(): Uint8Array {
    return new Uint8Array(0);
}

export function testUint8ArrayOutput_null(): Uint8Array | null {
    return null;
}

export function testUint8ClampedArrayInput(view: Uint8ClampedArray): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testUint8ClampedArrayOutput(): Uint8ClampedArray {
    const view = new Uint8ClampedArray(4);
    view.set([0, 1, 2, 3]);
    return view;
}

export function testUint16ArrayInput(view: Uint16Array): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testUint16ArrayOutput(): Uint16Array {
    const view = new Uint16Array(4);
    view.set([0, 1, 2, 3]);
    return view;
}

export function testUint32ArrayInput(view: Uint32Array): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testUint32ArrayOutput(): Uint32Array {
    const view = new Uint32Array(4);
    view.set([0, 1, 2, 3]);
    return view;
}

export function testUint64ArrayInput(view: Uint64Array): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testUint64ArrayOutput(): Uint64Array {
    const view = new Uint64Array(4);
    view.set([0, 1, 2, 3]);
    return view;
}

export function testFloat32ArrayInput(view: Float32Array): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testFloat32ArrayOutput(): Float32Array {
    const view = new Float32Array(4);
    view.set([0, 1, 2, 3]);
    return view;
}

export function testFloat64ArrayInput(view: Float64Array): void {
    assert(view.length == 4);
    assert(view[0] == 0);
    assert(view[1] == 1);
    assert(view[2] == 2);
    assert(view[3] == 3);
}

export function testFloat64ArrayOutput(): Float64Array {
    const view = new Float64Array(4);
    view.set([0, 1, 2, 3]);
    return view;
}
