/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// @ts-expect-error: decorator
@external("modus_system", "logMessage")
declare function hostLogMessage(level: string, message: string): void;

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
      hostLogMessage(traceLevel, message);
      break;
    case 1:
      hostLogMessage(traceLevel, `${message} ${a0}`);
      break;
    case 2:
      hostLogMessage(traceLevel, `${message} ${a0} ${a1}`);
      break;
    case 3:
      hostLogMessage(traceLevel, `${message} ${a0} ${a1} ${a2}`);
      break;
    case 4:
      hostLogMessage(traceLevel, `${message} ${a0} ${a1} ${a2} ${a3}`);
      break;
    case 5:
      hostLogMessage(traceLevel, `${message} ${a0} ${a1} ${a2} ${a3} ${a4}`);
      break;
    default:
      if (n < 0) hostLogMessage(traceLevel, message);
      else
        hostLogMessage(traceLevel, `${message} ${a0} ${a1} ${a2} ${a3} ${a4}`);
  }
}
