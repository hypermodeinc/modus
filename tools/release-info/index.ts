/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import {
  getLatestSdkVersion,
  getLatestRuntimeVersion,
  getLatestCliVersion,
  getAllSdkVersions,
  getAllRuntimeVersions,
  getAllCliVersions,
} from "./versioninfo.js";
import { SDK } from "./constants.js";
import { writeFileSync } from "fs";

async function buildModusLatestJSON(
  filename: string,
  includePrerelease: boolean,
) {
  const as = await getLatestSdkVersion(SDK.AssemblyScript, includePrerelease);
  const go = await getLatestSdkVersion(SDK.Go, includePrerelease);
  const cli = await getLatestCliVersion(includePrerelease);
  const runtime = await getLatestRuntimeVersion(includePrerelease);
  const info = {
    "sdk/assemblyscript": as,
    "sdk/go": go,
    cli: cli,
    runtime: runtime,
  };

  writeFileSync(filename, JSON.stringify(info, null, 2));
}

async function buildModusAllJSON(filename: string, includePrerelease: boolean) {
  const as = await getAllSdkVersions(SDK.AssemblyScript, includePrerelease);
  const go = await getAllSdkVersions(SDK.Go, includePrerelease);
  const cli = await getAllCliVersions(includePrerelease);
  const runtime = await getAllRuntimeVersions(includePrerelease);
  const info = {
    "sdk/assemblyscript": as,
    "sdk/go": go,
    cli: cli,
    runtime: runtime,
  };
  writeFileSync(filename, JSON.stringify(info, null, 2));
}

await buildModusLatestJSON("modus-latest.json", false);
await buildModusLatestJSON("modus-preview.json", true);
await buildModusAllJSON("modus-all.json", false);
await buildModusAllJSON("modus-preview-all.json", true);
