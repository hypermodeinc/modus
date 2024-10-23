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

function getHypEnvDir(): string {
  return path.join(process.env.HOME || "", ".hypermode");
}

function getSettingsFilePath(): string {
  return path.join(getHypEnvDir(), "settings.json");
}

export function readSettingsJson(): {
  email: null | string;
  jwt: null | string;
  orgId: null | string;
} {
  const content = fs.readFileSync(getSettingsFilePath(), "utf-8");

  let email: null | string = null;
  let jwt: null | string = null;
  let orgId: null | string = null;

  try {
    const jsonContent = JSON.parse(content);
    email = jsonContent.HYP_EMAIL || null;
    jwt = jsonContent.HYP_JWT || null;
    orgId = jsonContent.HYP_ORG_ID || null;
  } catch (e) {
    console.error("Error reading settings.json", e);
  }

  return { email, jwt, orgId };
}
