/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command, Flags } from "@oclif/core";
import chalk from "chalk";

import { execFileSync } from "node:child_process";
import { existsSync, mkdirSync, renameSync } from "node:fs";
import { rm } from "node:fs/promises";
import { arch, platform, tmpdir } from "node:os";
import path from "node:path";

import { clearLine, downloadFile, expandHomeDir } from "../../../util/index.js";
import { getLatestRuntimeVersion, runtimeReleaseExists, findCompatibleSdkVersion } from "../../../util/versioninfo.js";
import { GitHubOwner, GitHubRepo, SDK } from "../../../custom/globals.js";

export default class SDKInstallCommand extends Command {
  static args = {
    version: Args.string({
      description: "v0.0.0-|-SDK version to install",
      default: "latest",
    }),
  };

  static description = "Install a specific SDK version";
  static examples = ["modus sdk install v0.0.0", "modus sdk install latest"];

  static flags = {
    prerelease: Flags.boolean({
      char: "p",
      description: "Install a prerelease version",
    }),
  };

  async run(): Promise<void> {
    const { args, flags } = await this.parse(SDKInstallCommand);
    await this.installVersion(args.version, flags.prerelease);
  }

  async installVersion(version: string, prerelease: boolean) {
    if (version.toLowerCase() === "latest") {
      const versionText = prerelease ? "prerelease version" : "version";
      this.log(`[1/3] Getting latest ${versionText}`);
      const ver = await getLatestRuntimeVersion(prerelease);
      if (!version) {
        this.logError(`Failed to fetch latest ${versionText}`);
        this.exit(1);
      }
      version = ver!;
    } else {
      this.log(`[1/3] Checking version ${version}`);
      const exists = await runtimeReleaseExists(version!);
      if (!exists) {
        this.logError(`Version ${version} does not exist`);
        this.exit(1);
      }
    }

    clearLine();
    this.log(`[1/3] Found version ${version}`);

    this.log("[2/3] Downloading runtime ...");

    const tempDir = tmpdir();
    let osPlatform = platform().toString();
    let osArch = arch();
    if (osPlatform === "win32") osPlatform = "windows";
    if (osArch === "x64") osArch = "amd64";

    const archivePaths = [];

    const runtimeRelease = `runtime/${version}`;
    const runtimeFilename = `runtime_${version}_${osPlatform}_${osArch}.${osPlatform === "windows" ? "zip" : "tar.gz"}`;
    const runtimeDownloadUrl = `https://github.com/${GitHubOwner}/${GitHubRepo}/releases/download/${encodeURIComponent(runtimeRelease)}/${encodeURIComponent(runtimeFilename)}`;

    const runtimeArchivePath = path.join(tempDir, runtimeFilename);
    archivePaths.push({ archive: runtimeArchivePath, decompress: true });
    await downloadFile(runtimeDownloadUrl, runtimeArchivePath);

    // loop through each enum in SDK
    for (const sdk of Object.values(SDK)) {
      const sdkVersion = await findCompatibleSdkVersion(sdk, version);
      if (!sdkVersion) {
        this.logError(`Failed to find compatible ${sdk} SDK version for runtime ${version}`);
        this.exit(1);
      }

      clearLine();
      this.log(`[2/3] Downloading ${sdk} SDK version ${sdkVersion} ...`);

      const sdkRelease = `sdk/${sdk.toLowerCase()}/${sdkVersion}`;
      const sdkFilename = `templates_${sdk.toLowerCase()}_${sdkVersion}.tar.gz`;
      const sdkDownloadUrl = `https://github.com/${GitHubOwner}/${GitHubRepo}/releases/download/${encodeURIComponent(sdkRelease)}/${encodeURIComponent(sdkFilename)}`;

      const sdkArchivePath = path.join(tempDir, sdkFilename);
      archivePaths.push({ archive: sdkArchivePath, decompress: false });
      await downloadFile(sdkDownloadUrl, sdkArchivePath);
    }

    clearLine();
    this.log("[2/3] Downloads completed");

    this.log("[3/3] Installing ...");
    const installDir = expandHomeDir(`~/.modus/sdk/${version}/`);

    if (existsSync(installDir)) {
      await rm(installDir, { recursive: true, force: true });
    }

    for (const { archive, decompress } of archivePaths) {
      if (!existsSync(installDir)) {
        mkdirSync(installDir, { recursive: true });
      }
      if (decompress) {
        execFileSync("tar", ["-xf", archive, "-C", installDir]);
        await rm(archive);
      } else {
        renameSync(archive, path.join(installDir, path.basename(archive)));
      }
    }

    clearLine();
    this.log("[3/3] Installation successful");
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}
