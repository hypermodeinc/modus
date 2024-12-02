/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args } from "@oclif/core";
import chalk from "chalk";
import path from "node:path";
import os from "node:os";
import * as fs from "../../util/fs.js";
import * as vi from "../../util/versioninfo.js";
import { SDK } from "../../custom/globals.js";
import { getHeader } from "../../custom/header.js";
import { getAppInfo } from "../../util/appinfo.js";
import { withSpinner } from "../../util/index.js";
import { execFileWithExitCode } from "../../util/cp.js";
import SDKInstallCommand from "../sdk/install/index.js";
import { BaseCommand } from "../../baseCommand.js";

export default class BuildCommand extends BaseCommand {
  static args = {
    path: Args.directory({
      description: "Path to app directory",
      default: ".",
      exists: true,
    }),
  };

  static flags = {};

  static description = "Build a Modus app";

  static examples = ["modus build ./my-app"];

  async run(): Promise<void> {
    const { args, flags } = await this.parse(BuildCommand);

    const appPath = path.resolve(args.path);
    if (!(await fs.exists(path.join(appPath, "modus.json")))) {
      this.log(chalk.red("A modus.json file was not found at " + appPath));
      this.log(chalk.red("Please either execute the modus command from the app directory, or specify the path to the app you want to build."));
      this.exit(1);
    }

    const app = await getAppInfo(appPath);

    if (!flags["no-logo"]) {
      this.log(getHeader(this.config.version));
    }

    // pass chalk level to child processes so they can colorize output
    process.env.FORCE_COLOR = chalk.level.toString();

    if (!(await vi.sdkVersionIsInstalled(app.sdk, app.sdkVersion))) {
      await SDKInstallCommand.run([app.sdk, app.sdkVersion, "--no-logo"]);
    }

    const results = await withSpinner("Building " + app.name, async () => {
      const execOpts = {
        cwd: appPath,
        env: process.env,
        shell: true,
      };
      switch (app.sdk) {
        case SDK.AssemblyScript:
          if (!(await fs.exists(path.join(appPath, "node_modules")))) {
            const results = await execFileWithExitCode("npm", ["install"], execOpts);
            if (results.exitCode !== 0) {
              this.logError("Failed to install dependencies");
              return results;
            }
          }
          return await execFileWithExitCode("npx", ["modus-as-build"], execOpts);
        case SDK.Go: {
          const version = app.sdkVersion || (await vi.getLatestInstalledSdkVersion(app.sdk, true));
          if (!version) {
            this.logError("No installed version of the Modus Go SDK");
            return;
          }
          let buildTool = path.join(vi.getSdkPath(app.sdk, version), "modus-go-build");
          if (os.platform() === "win32") buildTool += ".exe";
          if (!(await fs.exists(buildTool))) {
            this.logError("Modus Go Build tool is not installed");
            return;
          }
          return await execFileWithExitCode(buildTool, ["."], execOpts);
        }
        default:
          this.logError("Unsupported SDK");
          this.exit(1);
      }
    });
    if (!results) {
      this.logError("No output from build command.");
      this.exit(1);
    } else if (results.exitCode !== 0) {
      this.log(chalk.bold.red("Build failed. ") + "😢\n");
      if (results.stdout) this.log(results.stdout.trimEnd() + "\n");
      if (results.stderr) this.logToStderr(results.stderr.trimEnd() + "\n");
      this.exit(results.exitCode);
    } else {
      this.log(chalk.bold.greenBright("Build succeeded! ") + "🎉\n");
      if (results.stdout) this.log(results.stdout.trimEnd() + "\n");
      if (results.stderr) this.logToStderr(results.stderr.trimEnd() + "\n");
    }
  }

  private logError(message: string) {
    this.log(chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}
