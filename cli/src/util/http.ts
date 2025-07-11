/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import ky from "ky";

// All outbound HTTP requests in the CLI should be made through functions in this file,
// to ensure consistent retry and timeout behavior.

export async function get<T>(url: string, timeout: number | false = 5000) {
  return await ky.get<T>(url, {
    retry: 4,
    timeout,
  });
}
