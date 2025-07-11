/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export class Metadata {
  // Installed runtimes - sorted by date
  static runtimes: string[] = [];

  static initialize(): void {
    // search hypermode.json / package.json for wanted runtime version
    // search current runtimes. populate
  }
}
