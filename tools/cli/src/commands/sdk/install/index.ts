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

import os from "node:os";
import path from "node:path";
import * as fs from "../../../util/fs.js";
import * as vi from "../../../util/versioninfo.js";
import * as globals from "../../../custom/globals.js";
import { execFile } from "../../../util/cp.js";
import { clearLine, downloadFile } from "../../../util/index.js";

export default class SDKInstallCommand extends Command {
  static args = {
    version: Args.string({
      description: "v0.0.0-|-SDK version to install",
      default: "latest",
    }),
  };

  static description = "Install a specific SDK version";
  static examples = ["modus sdk install v0.13.0", "modus sdk install", "modus sdk install latest", "modus sdk install latest --prerelease"];

  static flags = {
    force: Flags.boolean({
      char: "f",
      description: "Force re-installation if version already exists",
    }),
    prerelease: Flags.boolean({
      char: "p",
      aliases: ["pre"],
      description: "Install a prerelease version (used with 'latest' version)",
    }),
  };

  async run(): Promise<void> {
    const { args, flags } = await this.parse(SDKInstallCommand);
    await this.installVersion(args.version, flags.force, flags.prerelease);
  }

  async installVersion(version: string, force: boolean, prerelease: boolean) {
    if (version.toLowerCase() === "latest") {
      // first get the latest version that has been published
      const versionText = prerelease ? "prerelease version" : "version";
      this.log(`[1/3] Getting latest ${versionText}`);
      clearLine();
      const ver = await vi.getLatestRuntimeVersion(prerelease);
      if (!version) {
        this.logError(`Failed to fetch latest ${versionText}`);
        this.exit(1);
      }
      version = ver!;

      // then check if it's already installed
      if (await vi.versionIsInstalled(version)) {
        if (force) {
          this.log(chalk.yellow(version + " (latest) is already installed. Reinstalling ..."));
        } else {
          this.log(chalk.cyan(version + " (latest) is already installed."));
          this.exit(0);
        }
      }
    } else {
      // when specifying a version, _first_ check if it's already installed
      if (await vi.versionIsInstalled(version)) {
        if (force) {
          this.log(chalk.yellow(version + " is already installed. Reinstalling ..."));
        } else {
          this.log(chalk.cyan(version + " is already installed."));
          this.exit(0);
        }
      }

      // now check if the version exists online
      this.log(`[1/3] Checking version ${version}`);
      clearLine();
      const exists = await vi.runtimeReleaseExists(version!);
      if (!exists) {
        this.logError(`Version ${version} does not exist`);
        this.exit(1);
      }
    }

    this.log(`[1/3] Found version ${version}`);

    this.log("[2/3] Downloading runtime ...");

    const tempDir = os.tmpdir();
    let osPlatform = os.platform().toString();
    let osArch = os.arch();
    if (osPlatform === "win32") osPlatform = "windows";
    if (osArch === "x64") osArch = "amd64";

    const { GitHubOwner, GitHubRepo, ModusHomeDir } = globals;
    const archivePaths = [];

    const runtimeRelease = `runtime/${version}`;
    const runtimeFilename = `runtime_${version}_${osPlatform}_${osArch}.${osPlatform === "windows" ? "zip" : "tar.gz"}`;
    const runtimeDownloadUrl = `https://github.com/${GitHubOwner}/${GitHubRepo}/releases/download/${encodeURIComponent(runtimeRelease)}/${encodeURIComponent(runtimeFilename)}`;

    const runtimeArchivePath = path.join(tempDir, runtimeFilename);
    archivePaths.push({ archive: runtimeArchivePath, decompress: true });
    await downloadFile(runtimeDownloadUrl, runtimeArchivePath);

    for (const sdk of Object.values(globals.SDK)) {
      const sdkVersion = await vi.findCompatibleSdkVersion(sdk, version);
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
    const installDir = path.join(ModusHomeDir, "sdk", version);

    if (await fs.exists(installDir)) {
      await fs.rm(installDir, { recursive: true, force: true });
    }

    for (const { archive, decompress } of archivePaths) {
      if (!(await fs.exists(installDir))) {
        await fs.mkdir(installDir, { recursive: true });
      }
      if (decompress) {
        await execFile("tar", ["-xf", archive, "-C", installDir]);
        await fs.rm(archive);
      } else {
        await fs.rename(archive, path.join(installDir, path.basename(archive)));
      }
    }

    clearLine();
    this.log("[3/3] Installation successful");
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}
