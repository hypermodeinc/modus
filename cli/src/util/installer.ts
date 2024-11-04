/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import os from "node:os"
import path from "node:path"
import * as fs from "./fs.js"
import * as vi from "./versioninfo.js"
import { extract } from "./tar.js"
import { execFile } from "./cp.js"
import { downloadFile, isOnline } from "./index.js"
import { GitHubOwner, GitHubRepo, SDK } from "../custom/globals.js"

export async function installSDK(sdk: SDK, version: string) {
  if (!(await isOnline())) {
    throw new Error("No internet connection.  You must be online to install a Modus SDK.")
  }

  const sdkDir = vi.getSdkPath(sdk, version)
  const releaseTag = `sdk/${sdk.toLowerCase()}/${version}`
  const baseUrl = `https://github.com/${GitHubOwner}/${GitHubRepo}/releases/download/${encodeURIComponent(releaseTag)}/`

  const infoUrl = baseUrl + "sdk.json"
  if (!(await downloadFile(infoUrl, path.join(sdkDir, "sdk.json")))) {
    throw new Error(`Failed to download ${infoUrl}`)
  }

  const templatesUrl =
    baseUrl + encodeURIComponent(`templates_${sdk.toLowerCase()}_${version}.tar.gz`)
  if (!(await downloadFile(templatesUrl, path.join(sdkDir, "templates.tar.gz")))) {
    throw new Error(`Failed to download ${templatesUrl}`)
  }
}

export async function installRuntime(version: string) {
  if (!(await isOnline())) {
    throw new Error("No internet connection.  You must be online to install the Modus runtime.")
  }

  const tempDir = os.tmpdir()
  let osPlatform = os.platform().toString()
  let osArch = os.arch()
  if (osPlatform === "win32") osPlatform = "windows"
  if (osArch === "x64") osArch = "amd64"

  const releaseTag = `runtime/${version}`
  const fileName = `runtime_${version}_${osPlatform}_${osArch}.${osPlatform === "windows" ? "zip" : "tar.gz"}`
  const url = `https://github.com/${GitHubOwner}/${GitHubRepo}/releases/download/${encodeURIComponent(releaseTag)}/${encodeURIComponent(fileName)}`
  const archivePath = path.join(tempDir, fileName)

  if (!(await downloadFile(url, archivePath))) {
    throw new Error(`Failed to download ${url}`)
  }

  const installDir = vi.getRuntimePath(version)
  await fs.mkdir(installDir, { recursive: true })
  await extract(archivePath, installDir)
  await fs.rm(archivePath)
}

export async function installBuildTools(sdk: SDK, sdkVersion: string) {
  switch (sdk) {
    case SDK.Go:
      await installGoBuildTools(sdkVersion)
      break
  }
}

async function installGoBuildTools(sdkVersion: string) {
  const sdkPath = vi.getSdkPath(SDK.Go, sdkVersion)

  const ext = os.platform() === "win32" ? ".exe" : ""
  const buildTool = path.join(sdkPath, "modus-go-build" + ext)
  if (await fs.exists(buildTool)) {
    return
  }

  if (!(await isOnline())) {
    throw new Error("Could not find the Modus Go build tool. Please try again when you are online.")
  }

  const module = `github.com/${GitHubOwner}/${GitHubRepo}/sdk/go/tools/modus-go-build@${sdkVersion}`
  await execFile("go", ["install", module], {
    shell: true,
    env: {
      ...process.env,
      GOBIN: sdkPath,
    },
  })
}
