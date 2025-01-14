/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// @ts-expect-error: decorator
@external("modus_system", "logMessage")
declare function hostLogMessage(level: string, message: string): void;

// @ts-expect-error: decorator
@unsafe
@external("wasi_snapshot_preview1", "proc_exit")
declare function procExit(rval: u32): void;

export default function modus_abort(
  message: string | null = null,
  fileName: string | null = null,
  lineNumber: u32 = 0,
  columnNumber: u32 = 0,
): void {
  let msg = message ? message : "abort()";
  if (fileName) {
    msg += ` at ${fileName}:${lineNumber}:${columnNumber}`;
  }

  hostLogMessage("fatal", msg);
  procExit(255);
}
