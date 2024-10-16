/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

export const CLI_VERSION = "0.0.0";

export const GitHubOwner = "hypermodeinc";
export const GitHubRepo = "modus";
export const GitHubRuntimeTagPrefix = "runtime/";

// TODO: take this as a flag on CLI commands as needed
export const usePrereleaseVersions = true;

export enum SDK {
  AssemblyScript = "AssemblyScript",
  Go = "Go",
}
