/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { CLI_VERSION } from "../custom/globals.js";

export class Metadata {
  // CLI version (may want to have it read package.json instead)
  static cli_version: string = CLI_VERSION;
  // Runtime required by project/latest runtime
  static runtime_version: string = "0.12.6";
  // Installed runtimes - sorted by date
  static runtimes: string[] = [];
  static initialize(): void {
    // search hypermode.json / package.json for wanted runtime version
    // search current runtimes. populate
  }
  static async getLatestRuntime(): Promise<string | null> {
    try {
      // Its private for now. Need to update manually unfortunately
      // const response = await fetch("https://api.github.com/repos/HypermodeAI/runtime/releases/latest", {
      //     headers: {
      //         "Accept": "application/vnd.github.v3+json"
      //     }
      // });

      // if (!response.ok) {
      //     throw new Error(`Error fetching release: ${response.statusText}`);
      // }

      // const data = await response.json();
      // return data.tag_name || null; // Return the tag name if it exists
      return this.runtime_version;
    } catch {
      return null;
    }
  }
}
