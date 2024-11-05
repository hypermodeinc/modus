/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import fs from "node:fs";
import http from "node:http";
import net from "node:net";
import path from "node:path";
import { fileURLToPath } from "node:url";
import open from "open";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const filePath = path.resolve(__dirname, "../../content/index.html");

export async function openApiExplorer(endpoints: { [key: string]: string }) {
  const server = http.createServer((req, res) => {
    fs.readFile(filePath, (err, content) => {
      if (err) {
        res.writeHead(500);
        res.end(`Error loading file: ${err}`);
      } else {
        res.writeHead(200, { "Content-Type": "text/html" });
        res.end(content);
      }
    });
  });

  const port = await findAvailablePort(3000, 3100);
  server.listen(port, async () => {
    console.log(`API Explorer running at http://localhost:${port}`);
    await open(`http://localhost:${port}?endpoints=${encodeURIComponent(JSON.stringify(endpoints))}`);
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
