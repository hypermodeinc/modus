/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import os from "node:os";
import * as path from "node:path";
import * as fs from "./fs.js";
import chalk from "chalk";

type HypSettings = {
  email?: string;
  apiKey?: string;
  workspaceId?: string;
};

export async function readHypermodeSettings(): Promise<HypSettings> {
  const path = getSettingsFilePath();
  if (!(await fs.exists(path))) {
    return {};
  }

  try {
    const settings = JSON.parse(await fs.readFile(path, "utf-8"));
    return {
      email: settings.HYP_EMAIL,
      apiKey: settings.HYP_API_KEY,
      workspaceId: settings.HYP_WORKSPACE_ID,
    };
  } catch (e) {
    console.warn(chalk.yellow("Error reading " + path), e);
    return {};
  }
}

function getSettingsFilePath(): string {
  return path.join(os.homedir(), ".hypermode", "settings.json");
}
