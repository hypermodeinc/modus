/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { getLogo } from "./logo.js";

export function getHeader(cliVersion: string): string {
  let out = "";
  out += getLogo();
  out += "\n";
  out += `Modus CLI v${cliVersion}`;
  out += "\n";
  return out;
}
