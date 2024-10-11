/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// @ts-expect-error: decorator
@external("modus_system", "logMessage")
declare function log(level: string, message: string): void;

const traceLevel = "trace";

export default function modus_trace(
  message: string,
  n: i32 = 0,
  a0: f64 = 0,
  a1: f64 = 0,
  a2: f64 = 0,
  a3: f64 = 0,
  a4: f64 = 0,
): void {
  switch (n) {
    case 0:
      log(traceLevel, message);
      break;
    case 1:
      log(traceLevel, `${message} ${a0}`);
      break;
    case 2:
      log(traceLevel, `${message} ${a0} ${a1}`);
      break;
    case 3:
      log(traceLevel, `${message} ${a0} ${a1} ${a2}`);
      break;
    case 4:
      log(traceLevel, `${message} ${a0} ${a1} ${a2} ${a3}`);
      break;
    case 5:
      log(traceLevel, `${message} ${a0} ${a1} ${a2} ${a3} ${a4}`);
      break;
    default:
      if (n < 0) log(traceLevel, message);
      else log(traceLevel, `${message} ${a0} ${a1} ${a2} ${a3} ${a4}`);
  }
}
