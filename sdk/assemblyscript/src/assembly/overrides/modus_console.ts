/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// @ts-expect-error: decorator
@external("modus_system", "logMessage")
declare function hostLogMessage(level: string, message: string): void;

// @ts-expect-error: decorator
@lazy const timers = new Map<string, u64>();

export default abstract class modus_console {
  static assert<T>(condition: T, message: string = ""): void {
    if (!condition) {
      hostLogMessage("error", "Assertion failed: " + message);
    }
  }

  static log(message: string = ""): void {
    hostLogMessage("", message);
  }

  static debug(message: string = ""): void {
    hostLogMessage("debug", message);
  }

  static info(message: string = ""): void {
    hostLogMessage("info", message);
  }

  static warn(message: string = ""): void {
    hostLogMessage("warning", message);
  }

  static error(message: string = ""): void {
    hostLogMessage("error", message);
  }

  static time(label: string = "default"): void {
    const now = process.hrtime();
    if (timers.has(label)) {
      const msg = `Label '${label}' already exists for console.time()`;
      modus_console.warn(msg);
      return;
    }
    timers.set(label, now);
  }

  static timeLog(label: string = "default"): void {
    const now = process.hrtime();
    if (!timers.has(label)) {
      const msg = `No such label '${label}' for console.timeLog()`;
      modus_console.warn(msg);
      return;
    }
    modus_console.timeLogImpl(label, now);
  }

  static timeEnd(label: string = "default"): void {
    const now = process.hrtime();
    if (!timers.has(label)) {
      const msg = `No such label '${label}' for console.timeEnd()`;
      modus_console.warn(msg);
      return;
    }
    modus_console.timeLogImpl(label, now);
    timers.delete(label);
  }

  private static timeLogImpl(label: string, now: u64): void {
    const start = timers.get(label);
    const nanos = now - start;
    const millis = nanos / 1000000;
    modus_console.log(`${label}: ${millis}ms`);
  }
}
