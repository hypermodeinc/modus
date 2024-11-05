/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import * as fs from "../util/fs.js";
import net from "node:net";
import path from "node:path";
import { fileURLToPath } from "node:url";
import open from "open";
import express from "express";
import chalk from "chalk";

const MANIFEST_FILE = "modus.json";
const CONTENT_DIR = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../../content");

export async function openApiExplorer(appPath: string, runtimePort: number) {
  const app = express();
  app.use(express.static(CONTENT_DIR));

  app.get("/api/endpoints", async (_, res) => {
    const endpoints = await getGraphQLEndpointsFromManifest(appPath, runtimePort);
    res.json(endpoints);
  });

  const port = await findAvailablePort(3000, 3100);
  app.listen(port, async () => {
    const url = `http://localhost:${port}`;
    console.log(chalk.greenBright("Modus API Explorer is running at: ") + chalk.cyanBright(url) + "\n");
    await open(url);
  });
}

async function findAvailablePort(start: number, end: number): Promise<number> {
  for (let port = start; port <= end; port++) {
    const isAvailable = await checkPort(port);
    if (isAvailable) {
      return port;
    }
  }
  throw new Error(`No available port found in range ${start}-${end}`);
}

function checkPort(port: number): Promise<boolean> {
  return new Promise((resolve) => {
    const server = net.createServer();

    server.once("error", (err: NodeJS.ErrnoException) => {
      if (err.code === "EADDRINUSE") {
        resolve(false);
      } else {
        resolve(false);
      }
    });

    server.once("listening", () => {
      server.close(() => resolve(true));
    });

    server.listen(port);
  });
}

async function getGraphQLEndpointsFromManifest(appPath: string, port: number): Promise<{ [key: string]: string }> {
  const manifestPath = path.join(appPath, MANIFEST_FILE);
  if (!(await fs.exists(manifestPath))) {
    throw new Error(`Manifest file not found at ${manifestPath}`);
  }

  const manifestContent = await fs.readFile(manifestPath, "utf-8");
  const manifest = JSON.parse(manifestContent);

  if (!manifest.endpoints) {
    return {};
  }

  const results: { [key: string]: string } = {};
  for (const key in manifest.endpoints) {
    const ep = manifest.endpoints[key];
    results[key] = `http://localhost:${port}${ep.path}`;
  }
  return results;
}
