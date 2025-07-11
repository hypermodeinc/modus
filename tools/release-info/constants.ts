/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export const GitHubOwner = "hypermodeinc";
export const GitHubRepo = "modus";
export const GitHubRuntimeTagPrefix = "runtime/";
export const GitHubCliTagPrefix = "cli/";

export function GetSdkTagPrefix(sdk: SDK): string {
  return `sdk/${sdk.toLowerCase()}/`;
}

export enum SDK {
  AssemblyScript = "AssemblyScript",
  Go = "Go",
}
