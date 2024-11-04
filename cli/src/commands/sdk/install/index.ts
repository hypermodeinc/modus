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
import semver from "semver";

import * as fs from "../../../util/fs.js";
import * as vi from "../../../util/versioninfo.js";
import * as installer from "../../../util/installer.js";
import { isOnline, withSpinner } from "../../../util/index.js";
import { getHeader } from "../../../custom/header.js";
import { SDK, parseSDK } from "../../../custom/globals.js";
import { BaseCommand } from "../../../baseCommand.js";

export default class SDKInstallCommand extends BaseCommand {
  static args = {
    name: Args.string({
      description: "SDK name to install, or 'all' to install all SDKs",
      required: true,
    }),
    version: Args.string({
      description: "SDK version to install",
      default: "latest",
    }),
  };

  static description = "Install a Modus SDK";
  static examples = ["modus sdk install assemblyscript", "modus sdk install assemblyscript v1.2.3", "modus sdk install", "modus sdk install --prerelease"];

  static flags = {
    force: Flags.boolean({
      char: "f",
      default: false,
      description: "Force re-installation if version already exists",
    }),
    prerelease: Flags.boolean({
      char: "p",
      aliases: ["pre"],
      default: false,
      description: "Install a prerelease version (used with 'latest' version)",
    }),
  };

  async run(): Promise<void> {
    const { args, flags } = await this.parse(SDKInstallCommand);

    if (!flags["no-logo"]) {
      this.log(getHeader(this.config.version));
    }

    if (!(await isOnline())) {
      this.logError("No internet connection.  You must be online to install a Modus SDK.");
      this.exit(1);
    }

    const sdks: { sdk: SDK; version: string }[] = [];
    if (args.name.toLowerCase() === "all") {
      for (const sdk of Object.values(SDK)) {
        const version = await this.installSDK(sdk, args.version, flags.force, flags.prerelease);
        sdks.push({ sdk, version });
      }
    } else {
      const sdk = parseSDK(args.name);
      const version = await this.installSDK(sdk, args.version, flags.force, flags.prerelease);
      sdks.push({ sdk, version });
    }

    const runtimes = await this.getRuntimeVersions(sdks, flags.prerelease);
    for (const runtimeVersion of runtimes) {
      await this.installRuntime(runtimeVersion, flags.force);
    }

    this.log();
    this.log(chalk.bold.cyanBright("Installation successful!"));
    this.log();
  }

  private async installSDK(sdk: SDK, sdkVersion: string, force: boolean, prerelease: boolean): Promise<string> {
    let installSDK = true;
    let sdkText: string;
    if (sdkVersion.toLowerCase() === "latest") {
      sdkText = `Modus ${sdk} SDK${prerelease ? " prerelease" : ""}`;
      await withSpinner(chalk.dim(`Getting latest ${sdkText} version.`), async (spinner) => {
        const version = await vi.getLatestSdkVersion(sdk, prerelease);
        if (!version) {
          spinner.fail(chalk.red(`Failed to fetch latest ${sdkText} version.`));
          this.exit(1);
        }
        sdkVersion = version;
      });

      await withSpinner(chalk.dim(`Checking installations.`), async (spinner) => {
        sdkText = `Modus ${sdk} SDK ${sdkVersion}`;
        if (await vi.sdkVersionIsInstalled(sdk, sdkVersion)) {
          if (force) {
            spinner.warn(chalk.yellow(sdkText + " was already installed. Reinstalling."));
            const dir = vi.getSdkPath(sdk, sdkVersion);
            await fs.rm(dir, { recursive: true, force: true });
          } else {
            spinner.succeed(chalk.dim(sdkText + " was already installed."));
            installSDK = false;
          }
        }
      });
    } else if (!sdkVersion.startsWith("v")) {
      this.logError("Version must start with 'v'.");
      this.exit(1);
    } else {
      sdkText = `Modus ${sdk} SDK ${sdkVersion}`;

      await withSpinner(chalk.dim(`Checking installations.`), async (spinner) => {
        if (await vi.sdkVersionIsInstalled(sdk, sdkVersion)) {
          if (force) {
            spinner.warn(chalk.yellow(sdkText + " was already installed. Reinstalling."));
          } else {
            spinner.succeed(chalk.dim(sdkText + " was already installed."));
            installSDK = false;
          }
        }
      });

      await withSpinner(chalk.dim(`Checking version ${sdkVersion}`), async (spinner) => {
        const exists = await vi.sdkReleaseExists(sdk, sdkVersion);
        if (!exists) {
          spinner.fail(chalk.red(sdkText + " does not exist!"));
          this.exit(1);
        }
      });
    }

    if (installSDK) {
      await withSpinner(chalk.dim("Downloading and installing " + sdkText), async (spinner) => {
        try {
          await installer.installSDK(sdk, sdkVersion);
        } catch (e) {
          spinner.fail(chalk.red(`Failed to download ${sdkText}`));
          throw e;
        }

        await installer.installBuildTools(sdk, sdkVersion);

        spinner.succeed(chalk.dim(`Installed ${sdkText}`));
      });
    }

    return sdkVersion;
  }

  private async getRuntimeVersions(sdks: { sdk: SDK; version: string }[], prerelease: boolean): Promise<string[]> {
    const runtimes = new Set<string>();
    for (const { sdk, version } of sdks) {
      const sdkText = `Modus ${sdk} SDK ${version}`;
      await withSpinner(chalk.dim("Checking for latest runtime compatible with " + sdkText), async (spinner) => {
        const runtimeVersion = await vi.findLatestCompatibleRuntimeVersion(sdk, version, prerelease);
        if (!runtimeVersion) {
          spinner.fail(chalk.red("Failed to find compatible runtime for " + sdkText));
          this.exit(1);
        }
        runtimes.add(runtimeVersion);
      });
    }

    return semver.rsort(Array.from(runtimes).map((v) => v.slice(1))).map((v) => "v" + v);
  }

  private async installRuntime(version: string, force: boolean): Promise<void> {
    const runtimeText = `Modus Runtime ${version}`;
    const installRuntime = await withSpinner(chalk.dim("Checking for runtime installation."), async (spinner) => {
      if (await vi.runtimeVersionIsInstalled(version)) {
        if (force) {
          spinner.warn(chalk.yellow(runtimeText + " was already installed. Reinstalling."));
          const dir = vi.getRuntimePath(version);
          await fs.rm(dir, { recursive: true, force: true });
          return true;
        } else {
          spinner.succeed(chalk.dim(runtimeText + " was already installed."));
          return false;
        }
      }
      return true;
    });

    if (installRuntime) {
      await withSpinner(chalk.dim("Downloading and installing " + runtimeText), async (spinner) => {
        try {
          await installer.installRuntime(version);
        } catch (e) {
          spinner.fail(chalk.red(`Failed to download ${runtimeText}`));
          throw e;
        }
        spinner.succeed(chalk.dim(`Installed ${runtimeText}`));
      });
    }
  }

  private logError(message: string) {
    this.log(chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}
