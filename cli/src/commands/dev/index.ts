/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command, Flags } from "@oclif/core";
import { spawn } from "node:child_process";
import path from "node:path";
import os from "node:os";
import chalk from "chalk";
import pm from "picomatch";
import chokidar from "chokidar";

import * as fs from "../../util/fs.js";
import * as vi from "../../util/versioninfo.js";
import * as installer from "../../util/installer.js";
import { SDK } from "../../custom/globals.js";
import { getHeader } from "../../custom/header.js";
import { getAppInfo } from "../../util/appinfo.js";
import { isOnline, withSpinner } from "../../util/index.js";
import { readHypermodeSettings } from "../../util/hypermode.js";
import BuildCommand from "../build/index.js";

export default class DevCommand extends Command {
  static args = {
    path: Args.directory({
      description: "Path to app directory",
      default: ".",
      exists: true,
    }),
  };

  static flags = {
    help: Flags.help({
      char: "h",
      helpLabel: "-h, --help",
      description: "Show help message",
    }),
    "no-logo": Flags.boolean({
      aliases: ["no-logo"],
      hidden: true,
    }),
    runtime: Flags.string({
      char: "r",
      description: "Modus runtime version to use. If not provided, the latest runtime compatible with the app will be used.",
    }),
    prerelease: Flags.boolean({
      char: "p",
      aliases: ["pre"],
      description: "Use a prerelease version of the Modus runtime. Not needed if specifying a runtime version.",
    }),
    "no-build": Flags.boolean({
      aliases: ["nobuild"],
      description: "Don't build the app before running (implies --no-watch)",
    }),
    "no-watch": Flags.boolean({
      aliases: ["nowatch"],
      description: "Don't watch app code for changes",
    }),
    delay: Flags.integer({
      description: "Delay (in milliseconds) between file change detection and rebuild",
      default: 500,
    }),
  };

  static description = "Run a Modus app locally for development";

  static examples = ["modus dev", "modus dev ./my-app", "modus dev ./my-app --no-watch"];

  async run(): Promise<void> {
    const { args, flags } = await this.parse(DevCommand);

    const appPath = path.resolve(args.path);
    if (!(await fs.exists(path.join(appPath, "modus.json")))) {
      this.log(chalk.red("A modus.json file was not found at " + appPath));
      this.log(chalk.red("Please either execute the modus command from the app directory, or specify the path to the app you want to run."));
      this.exit(1);
    }

    const app = await getAppInfo(appPath);
    const { sdk, sdkVersion } = app;

    if (!flags["no-logo"]) {
      this.log(getHeader(this.config.version));
    }

    if (!(await vi.sdkVersionIsInstalled(sdk, sdkVersion))) {
      const sdkText = `Modus ${sdk} SDK ${sdkVersion}`;
      await withSpinner(chalk.dim("Downloading and installing " + sdkText), async (spinner) => {
        try {
          await installer.installSDK(sdk, sdkVersion);
        } catch (e) {
          spinner.fail(chalk.red(`Failed to download ${sdkText}`));
          throw e;
        }
        spinner.succeed(chalk.dim(`Installed ${sdkText}`));
      });
    }

    let runtimeVersion = flags.runtime;
    if (runtimeVersion) {
      const runtimeText = `Modus Runtime ${runtimeVersion}`;
      if (!(await vi.runtimeVersionIsInstalled(runtimeVersion))) {
        if (await isOnline()) {
          await withSpinner(chalk.dim("Downloading and installing " + runtimeText), async (spinner) => {
            try {
              await installer.installRuntime(runtimeVersion!);
            } catch (e) {
              spinner.fail(chalk.red("Failed to download " + runtimeText));
              throw e;
            }
            spinner.succeed(chalk.dim("Installed " + runtimeText));
          });
        } else {
          this.logError(`${runtimeText} is not installed, and you are offline. Please try again when you have an internet connection.`);
          this.exit(1);
        }
      }
    } else if (await isOnline()) {
      const version = await vi.findLatestCompatibleRuntimeVersion(sdk, sdkVersion, flags.prerelease);
      if (version && !(await vi.runtimeVersionIsInstalled(version))) {
        const runtimeText = `Modus Runtime ${version}`;
        await withSpinner(chalk.dim("Downloading and installing " + runtimeText), async (spinner) => {
          try {
            await installer.installRuntime(version!);
          } catch (e) {
            spinner.fail(chalk.red("Failed to download " + runtimeText));
            throw e;
          }
          spinner.succeed(chalk.dim("Installed " + runtimeText));
        });
      }
      if (!version) {
        this.logError("Could not find a compatible Modus runtime version. Please try again.");
        return;
      }
      runtimeVersion = version;
    } else {
      const version = await vi.findCompatibleInstalledRuntimeVersion(sdk, sdkVersion, flags.prerelease);
      if (!version) {
        this.logError("Could not find a compatible Modus runtime version. Please try again when you have an internet connection.");
        return;
      }
      runtimeVersion = version;
    }

    const ext = os.platform() === "win32" ? ".exe" : "";
    const runtimePath = path.join(vi.getRuntimePath(runtimeVersion), "modus_runtime" + ext);

    if (!flags["no-build"]) {
      await BuildCommand.run([appPath, "--no-logo"]);
    }

    const hypSettings = await readHypermodeSettings();

    const env = {
      ...process.env,
      MODUS_ENV: "dev",
      HYP_EMAIL: hypSettings.email,
      HYP_JWT: hypSettings.jwt,
      HYP_ORG_ID: hypSettings.orgId,
    };

    const child = spawn(runtimePath, ["-appPath", path.join(appPath, "build")], {
      stdio: ["inherit", "inherit", "pipe"],
      env: env,
    });
    child.stderr.pipe(process.stderr);
    child.on("close", (code) => this.exit(code || 1));

    if (!flags["no-watch"] && !flags["no-build"]) {
      let lastModified = 0;
      let lastBuild = 0;
      let paused = true;
      setInterval(async () => {
        if (paused) {
          return;
        }
        paused = true;

        if (lastBuild > lastModified) {
          return;
        }
        lastBuild = Date.now();

        try {
          child.stderr.pause();
          this.log();
          this.log(chalk.magentaBright("Detected change. Rebuilding..."));
          this.log();
          await BuildCommand.run([appPath, "--no-logo"]);
        } catch (e) {
          this.log(chalk.magenta("Waiting for more changes..."));
          this.log(chalk.dim("Press Ctrl+C at any time to stop the server."));
        } finally {
          child.stderr.resume();
        }
      }, flags.delay);

      const globs = getGlobsToWatch(sdk);

      // NOTE: The built-in fs.watch or fsPromises.watch is insufficient for our needs.
      // Instead, we use chokidar for consistent behavior in cross-platform file watching.
      const pmOpts: pm.PicomatchOptions = { posixSlashes: true };
      chokidar
        .watch(appPath, {
          ignored: (filePath, stats) => {
            const relativePath = path.relative(appPath, filePath);
            if (!stats || !relativePath) return false;

            let ignore = false;
            if (pm(globs.excluded, pmOpts)(relativePath)) {
              ignore = true;
            } else if (stats.isFile()) {
              ignore = !pm(globs.included, pmOpts)(relativePath);
            }

            if (process.env.MODUS_DEBUG) {
              this.log(chalk.dim(`${ignore ? "ignored: " : "watching:"}  ${relativePath}`));
            }
            return ignore;
          },
          ignoreInitial: true,
          persistent: true,
        })
        .on("all", () => {
          lastModified = Date.now();
          paused = false;
        });
    }
  }

  private logError(message: string) {
    this.log(chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

function getGlobsToWatch(sdk: SDK) {
  const included: string[] = [];
  const excluded: string[] = [".git/**", "build/**"];

  switch (sdk) {
    case SDK.AssemblyScript:
      included.push("**/*.ts", "**/asconfig.json", "**/tsconfig.json", "**/package.json");
      excluded.push("node_modules/**");
      break;

    case SDK.Go:
      included.push("**/*.go", "**/go.mod");
      excluded.push("**/*_generated.go", "**/*.generated.go", "**/*_test.go");
      break;

    default:
      throw new Error(`Unsupported SDK: ${sdk}`);
  }
  return { included, excluded };
}
