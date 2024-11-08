/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import ky from "ky";

// All outbound HTTP requests in the CLI should be made through functions in this file,
// to ensure consistent retry and timeout behavior.

export async function get<T>(url: string, timeout: number | false = 10000) {
  return await ky.get<T>(url, {
    retry: 10, // retry 10 times, up to the default 10s timeout (unless overridden)
    timeout,
  });
}
