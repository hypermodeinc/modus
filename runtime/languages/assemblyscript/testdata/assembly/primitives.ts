/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export function testBoolInput_false(b: bool): void {
  assert(b == false);
}

export function testBoolInput_true(b: bool): void {
  assert(b == true);
}

export function testBoolOutput_false(): bool {
  return false;
}

export function testBoolOutput_true(): bool {
  return true;
}

export function testI8Input_min(n: i8): void {
  assert(n == i8.MIN_VALUE);
}

export function testI8Input_max(n: i8): void {
  assert(n == i8.MAX_VALUE);
}

export function testI8Output_min(): i8 {
  return i8.MIN_VALUE;
}

export function testI8Output_max(): i8 {
  return i8.MAX_VALUE;
}

export function testI16Input_min(n: i16): void {
  assert(n == i16.MIN_VALUE);
}

export function testI16Input_max(n: i16): void {
  assert(n == i16.MAX_VALUE);
}

export function testI16Output_min(): i16 {
  return i16.MIN_VALUE;
}

export function testI16Output_max(): i16 {
  return i16.MAX_VALUE;
}

export function testI32Input_min(n: i32): void {
  assert(n == i32.MIN_VALUE);
}

export function testI32Input_max(n: i32): void {
  assert(n == i32.MAX_VALUE);
}

export function testI32Output_min(): i32 {
  return i32.MIN_VALUE;
}

export function testI32Output_max(): i32 {
  return i32.MAX_VALUE;
}

export function testI64Input_min(n: i64): void {
  assert(n == i64.MIN_VALUE);
}

export function testI64Input_max(n: i64): void {
  assert(n == i64.MAX_VALUE);
}

export function testI64Output_min(): i64 {
  return i64.MIN_VALUE;
}

export function testI64Output_max(): i64 {
  return i64.MAX_VALUE;
}

export function testISizeInput_min(n: isize): void {
  assert(n == isize.MIN_VALUE);
}

export function testISizeInput_max(n: isize): void {
  assert(n == isize.MAX_VALUE);
}

export function testISizeOutput_min(): isize {
  return isize.MIN_VALUE;
}

export function testISizeOutput_max(): isize {
  return isize.MAX_VALUE;
}

export function testU8Input_min(n: u8): void {
  assert(n == u8.MIN_VALUE);
}

export function testU8Input_max(n: u8): void {
  assert(n == u8.MAX_VALUE);
}

export function testU8Output_min(): u8 {
  return u8.MIN_VALUE;
}

export function testU8Output_max(): u8 {
  return u8.MAX_VALUE;
}

export function testU16Input_min(n: u16): void {
  assert(n == u16.MIN_VALUE);
}

export function testU16Input_max(n: u16): void {
  assert(n == u16.MAX_VALUE);
}

export function testU16Output_min(): u16 {
  return u16.MIN_VALUE;
}

export function testU16Output_max(): u16 {
  return u16.MAX_VALUE;
}

export function testU32Input_min(n: u32): void {
  assert(n == u32.MIN_VALUE);
}

export function testU32Input_max(n: u32): void {
  assert(n == u32.MAX_VALUE);
}

export function testU32Output_min(): u32 {
  return u32.MIN_VALUE;
}

export function testU32Output_max(): u32 {
  return u32.MAX_VALUE;
}

export function testU64Input_min(n: u64): void {
  assert(n == u64.MIN_VALUE);
}

export function testU64Input_max(n: u64): void {
  assert(n == u64.MAX_VALUE);
}

export function testU64Output_min(): u64 {
  return u64.MIN_VALUE;
}

export function testU64Output_max(): u64 {
  return u64.MAX_VALUE;
}

export function testUSizeInput_min(n: usize): void {
  assert(n == usize.MIN_VALUE);
}

export function testUSizeInput_max(n: usize): void {
  assert(n == usize.MAX_VALUE);
}

export function testUSizeOutput_min(): usize {
  return usize.MIN_VALUE;
}

export function testUSizeOutput_max(): usize {
  return usize.MAX_VALUE;
}

export function testF32Input_min(n: f32): void {
  assert(n == f32.MIN_VALUE);
}

export function testF32Input_max(n: f32): void {
  assert(n == f32.MAX_VALUE);
}

export function testF32Output_min(): f32 {
  return f32.MIN_VALUE;
}

export function testF32Output_max(): f32 {
  return f32.MAX_VALUE;
}

export function testF64Input_min(n: f64): void {
  assert(n == f64.MIN_VALUE);
}

export function testF64Input_max(n: f64): void {
  assert(n == f64.MAX_VALUE);
}

export function testF64Output_min(): f64 {
  return f64.MIN_VALUE;
}

export function testF64Output_max(): f64 {
  return f64.MAX_VALUE;
}
