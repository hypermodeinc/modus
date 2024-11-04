/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
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
