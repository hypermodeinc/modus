/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// "Hello World" in Japanese
const testString = "こんにちは、世界";

export function testStringInput(s: string): void {
  assert(s == testString);
}

export function testStringOutput(): string {
  return testString;
}

export function testStringInput_empty(s: string): void {
  assert(s == "");
}

export function testStringOutput_empty(): string {
  return "";
}

export function testNullStringInput(s: string | null): void {
  assert(s == testString);
}

export function testNullStringOutput(): string | null {
  return testString;
}

export function testNullStringInput_empty(s: string | null): void {
  assert(s == "");
}

export function testNullStringOutput_empty(): string | null {
  return "";
}

export function testNullStringInput_null(s: string | null): void {
  assert(s == null);
}

export function testNullStringOutput_null(): string | null {
  return null;
}
