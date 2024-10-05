/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// @ts-expect-error: decorator
@external("hypermode", "log")
declare function log(level: string, message: string): void;

// @ts-expect-error: decorator
@lazy const timers = new Map<string, u64>();

export default abstract class modus_console {
  static assert<T>(condition: T, message: string = ""): void {
    if (!condition) {
      log("error", "Assertion failed: " + message);
    }
  }

  static log(message: string = ""): void {
    log("", message);
  }

  static debug(message: string = ""): void {
    log("debug", message);
  }

  static info(message: string = ""): void {
    log("info", message);
  }

  static warn(message: string = ""): void {
    log("warning", message);
  }

  static error(message: string = ""): void {
    log("error", message);
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
