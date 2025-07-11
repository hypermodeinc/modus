/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import path from "node:path";
import os from "node:os";
import process from "node:process";

export const ModusHomeDir = process.env.MODUS_HOME || path.join(os.homedir(), ".modus");

export const MinNodeVersion = "22.0.0";
export const MinGoVersion = "1.23.1";
export const MinTinyGoVersion = "0.35.0";

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
