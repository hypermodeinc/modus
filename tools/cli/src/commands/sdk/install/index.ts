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
import { existsSync, mkdirSync } from "node:fs";
import { rm } from "node:fs/promises";
import { arch, platform, tmpdir } from "node:os";
import path from "node:path";

import { clearLine, downloadFile, expandHomeDir } from "../../../util/index.js";
import { getLatestRuntimeVersion, runtimeReleaseExists } from "../../../util/versioninfo.js";
import { GitHubOwner, GitHubRepo } from "../../../custom/globals.js";

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

    this.log("[2/3] Downloading ...");

    let osPlatform = platform().toString();
    let osArch = arch();
    if (osPlatform === "win32") osPlatform = "windows";
    if (osArch === "x64") osArch = "amd64";

    const release = `runtime/${version}`;
    const filename = `runtime_${version}_${osPlatform}_${osArch}.${osPlatform === "windows" ? "zip" : "tar.gz"}`;
    const downloadUrl = `https://github.com/${GitHubOwner}/${GitHubRepo}/releases/download/${encodeURIComponent(release)}/${encodeURIComponent(filename)}`;

    const archiveName = "modus-" + filename;
    const archivePath = path.join(tmpdir(), archiveName);
    await downloadFile(downloadUrl, archivePath);

    clearLine();
    this.log("[2/3] Download completed");

    this.log("[3/3] Installing ...");
    const installDir = expandHomeDir(`~/.modus/sdk/${version}/`);

    if (existsSync(installDir)) {
      await rm(installDir, { recursive: true, force: true });
    }

    mkdirSync(installDir, { recursive: true });
    execFileSync("tar", ["-xf", archivePath, "-C", installDir]);
    await rm(archivePath);

    clearLine();
    this.log("[3/3] Installation successful");
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}
