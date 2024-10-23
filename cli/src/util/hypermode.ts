/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import * as path from "node:path";
import * as fs from "node:fs";
import * as utils from "./fs.js";
import chalk from "chalk";

function getHypEnvDir(): string {
  return path.join(process.env.HOME || "", ".hypermode");
}

function getSettingsFilePath(): string {
  return path.join(getHypEnvDir(), "settings.json");
}

type HypSettings = {
  email?: string;
  jwt?: string;
  orgId?: string;
};

export async function readSettingsJson(): Promise<HypSettings> {
  const path = getSettingsFilePath();
  const content = await utils.readFile(path, "utf-8");

  const settings: HypSettings = {};

  try {
    const jsonContent = JSON.parse(content);
    settings.email = jsonContent.HYP_EMAIL;
    settings.jwt = jsonContent.HYP_JWT;
    settings.orgId = jsonContent.HYP_ORG_ID;
  } catch (e) {
    console.warn(chalk.yellow("Error reading " + path), e);
  }

  return settings;
}
