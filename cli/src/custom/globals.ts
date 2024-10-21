/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import path from "node:path";
import os from "node:os";

export const ModusHomeDir = path.join(os.homedir(), ".modus");

export const MinNodeVersion = "22.0.0";
export const MinGoVersion = "1.23.0";
export const MinTinyGoVersion = "0.33.0";

export const GitHubOwner = "hypermodeinc";
export const GitHubRepo = "modus";
export const GitHubRuntimeTagPrefix = "runtime/";

export function GetSdkTagPrefix(sdk: SDK): string {
  return `sdk/${sdk.toLowerCase()}/`;
}

export enum SDK {
  AssemblyScript = "AssemblyScript",
  Go = "Go",
}

export function parseSDK(sdk: string): SDK {
  switch (sdk.toLowerCase()) {
    case "as":
    case "assemblyscript":
      return SDK.AssemblyScript;
    case "go":
    case "golang":
      return SDK.Go;
    default:
      throw new Error(`Unknown SDK: ${sdk}`);
  }
}
