/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export function testArrayBufferInput(buffer: ArrayBuffer): void {
  const view = Uint8Array.wrap(buffer);
  assert(view.length == 4);
  assert(view[0] == 1);
  assert(view[1] == 2);
  assert(view[2] == 3);
  assert(view[3] == 4);
}

export function testArrayBufferOutput(): ArrayBuffer {
  const buffer = new ArrayBuffer(4);
  const view = Uint8Array.wrap(buffer);
  view[0] = 1;
  view[1] = 2;
  view[2] = 3;
  view[3] = 4;
  return buffer;
}
