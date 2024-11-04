/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { getLogo } from "./logo.js"

export function getHeader(cliVersion: string): string {
  let out = ""
  out += getLogo()
  out += "\n"
  out += `Modus CLI v${cliVersion}`
  out += "\n"
  return out
}
