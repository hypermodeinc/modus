/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command } from "@oclif/core";
import chalk from "chalk";
import path from "node:path";
import os from "node:os";
import { execFileSync } from "node:child_process";
import * as fs from "../../util/fs.js";
import * as vi from "../../util/versioninfo.js";
import { ModusHomeDir, SDK } from "../../custom/globals.js";

export default class BuildCommand extends Command {
  static args = {
    path: Args.string({
      description: "./my-project-|-Directory to build",
      default: ".",
    }),
  };

  static description = "Build a Modus project";

  static examples = ["modus build ./my-project"];

  static flags = {};

  async run(): Promise<void> {
    const { args } = await this.parse(BuildCommand);

    const cwd = path.resolve(args.path);
    const sdk = await determineSDK(cwd);

    // TODO: Reuse the logic from the `modus new` command to check for latest versions when online.
    // But also, we need to somehow figure out which version of the runtime is compatible with the project.
    // For now, we'll just use the latest installed version.
    const installedVersion = await vi.getLatestInstalledVersion();
    if (!installedVersion) {
      this.logError("No Modus SDK installed");
      return;
    }

    switch (sdk) {
      case SDK.AssemblyScript:
        if (!(await fs.exists(path.join(cwd, "node_modules")))) {
          execFileSync("npm", ["install"], { cwd, env: process.env, stdio: "inherit" });
        }
        execFileSync("npx", ["modus-as-build"], { cwd, env: process.env, stdio: "inherit" });
        break;

      case SDK.Go:
        let buildTool = path.join(ModusHomeDir, "sdk", installedVersion, "modus-go-build");
        if (os.platform() === "win32") buildTool += ".exe";
        if (!(await fs.exists(buildTool))) {
          this.logError("Modus Go Build tool is not installed");
          return;
        }

        execFileSync(buildTool, ["."], { cwd, env: process.env, stdio: "inherit" });

        break;
    }
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

async function determineSDK(appPath: string): Promise<SDK> {
  if (await fs.exists(path.join(appPath, "package.json"))) {
    return SDK.AssemblyScript;
  }
  if (await fs.exists(path.join(appPath, "go.mod"))) {
    return SDK.Go;
  }
  throw new Error("Could not determine SDK");
}
