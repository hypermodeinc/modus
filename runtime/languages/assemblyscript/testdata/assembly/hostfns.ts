/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// @ts-expect-error: decorator
@external("modus_test", "add")
declare function hostAdd(a: i32, b: i32): i32

// @ts-expect-error: decorator
@external("modus_test", "echo")
declare function hostEcho(message: string): string

// @ts-expect-error: decorator
@external("modus_test", "echoObject")
declare function hostEchoObject(obj: TestHostObject): TestHostObject

class TestHostObject {
  a!: i32
  b!: bool
  c!: string
}

export function add(a: i32, b: i32): i32 {
  return hostAdd(a, b)
}

export function echo(message: string): string {
  return hostEcho(message)
}

export function echoObject(obj: TestHostObject): TestHostObject {
  return hostEchoObject(obj)
}
