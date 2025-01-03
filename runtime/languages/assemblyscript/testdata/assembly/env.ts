/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export function now(): i64 {
  return Date.now()
}

export function spin(duration: i64): i64 {
  const start = performance.now()

  let d = Date.now()
  while (Date.now() - d <= duration) {
    // do nothing
  }

  const end = performance.now()

  return i64(end - start)
}
