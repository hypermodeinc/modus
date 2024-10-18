/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Command, Flags } from "@oclif/core";
import chalk from "chalk";

import os from "node:os";
import path from "node:path";

import * as fs from "../../util/fs.js";
import * as vi from "../../util/versioninfo.js";
import { execFile } from "../../util/cp.js";
import { isOnline } from "../../util/index.js";
import { GitHubOwner, GitHubRepo, ModusHomeDir, SDK, parseSDK } from "../../custom/globals.js";
import { ask, clearLine, withSpinner } from "../../util/index.js";
import SDKInstallCommand from "../sdk/install/index.js";

export default class NewCommand extends Command {
  static description = "Create a new Modus app";

  static examples = ["modus new", "modus new --name my-app", "modus new --name my-app --sdk go --dir ./my-app --force"];

  static flags = {
    name: Flags.string({
      char: "n",
      description: "App name",
    }),
    dir: Flags.string({
      char: "d",
      aliases: ["directory"],
      description: "App directory",
    }),
    sdk: Flags.string({
      char: "s",
      description: "SDK to use",
    }),
    template: Flags.string({
      char: "t",
      description: "Template to use",
    }),
    force: Flags.boolean({
      char: "f",
      default: false,
      description: "Initialize without prompting",
    }),
    prerelease: Flags.boolean({
      char: "p",
      aliases: ["pre"],
      default: false,
      description: "Use a prerelease version of the Modus SDK",
    }),
  };

  async run(): Promise<void> {
    const { flags } = await this.parse(NewCommand);

    this.log(chalk.bold(`Modus CLI v${this.config.version}`));

    const name = flags.name || (await this.promptAppName());
    if (!name) {
      this.logError("An app name is required.");
      this.exit(1);
    }

    const dir = flags.dir ? path.join(process.cwd(), flags.dir) : await this.promptInstallPath("." + path.sep + name);
    if (!dir) {
      this.logError("An install directory is required.");
      this.exit(1);
    }

    const sdk = parseSDK(
      flags.sdk
        ? Object.values(SDK)[
            Object.keys(SDK)
              .map((v) => v.toLowerCase())
              .indexOf(flags.sdk?.trim().toLowerCase())
          ]
        : await this.promptSdkSelection()
    );

    const template = flags.template || (await this.promptTemplate("default"));
    if (!template) {
      this.logError("A template is required.");
      this.exit(1);
    }

    if (!flags.force && !(await this.confirmAction("[5/5] Continue? [y/n]"))) {
      this.log(chalk.dim("Aborted."));
      this.exit(1);
    }

    await this.createApp(name, dir, sdk, template, flags.force, flags.prerelease);
  }

  private async promptAppName(): Promise<string> {
    this.log("[1/5] App Name:");
    const name = ((await ask(chalk.dim(" -> "))) || "").trim();
    clearLine();
    clearLine();
    this.log("[1/5] App Name: " + chalk.dim(name.length ? name : "Not Provided"));
    return name;
  }

  private async promptInstallPath(defaultValue: string): Promise<string> {
    this.log("[2/5] Install Dir: " + chalk.dim(`(${defaultValue})`));
    const dir = ((await ask(chalk.dim(" -> "))) || defaultValue).trim();
    clearLine();
    clearLine();
    this.log("[2/5] Directory: " + chalk.dim(dir));
    return path.resolve(dir);
  }

  private async promptSdkSelection(): Promise<string> {
    this.log("[3/5] Select an SDK");
    for (const [index, sdk] of Object.values(SDK).entries()) {
      this.log(chalk.dim(` ${index + 1}. ${sdk}`));
    }

    const selectedIndex = Number.parseInt(((await ask(chalk.dim(" -> "))) || "1").trim(), 10) - 1;
    const sdk = Object.values(SDK)[selectedIndex];
    clearLine();
    clearLine();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    for (const _ of Object.values(SDK)) clearLine();
    if (!sdk) this.exit(1);
    this.log("[3/5] SDK: " + chalk.dim(sdk));
    return sdk;
  }

  private async promptTemplate(defaultValue: string): Promise<string> {
    this.log("[4/5] Template: " + chalk.dim(`(${defaultValue})`));
    const template = ((await ask(chalk.dim(" -> "))) || defaultValue).trim();
    clearLine();
    clearLine();
    this.log("[4/5] Template: " + chalk.dim(template));
    return template;
  }

  private async createApp(name: string, dir: string, sdk: SDK, template: string, force: boolean, prerelease: boolean) {
    if (!force && (await fs.exists(dir))) {
      if (!(await this.confirmAction("Attempting to overwrite a folder that already exists.\nAre you sure you want to continue? [y/n]"))) {
        clearLine();
        return;
      } else {
        clearLine();
        clearLine();
      }
    }

    // Validate SDK-specific prerequisites
    switch (sdk) {
      case SDK.AssemblyScript:
        if (parseInt(process.versions.node.split(".")[0]) < 22) {
          this.logError("The Modus AssemblyScript SDK requires Node.js v22 or later.");
          this.exit(1);
        }
        break;
      case SDK.Go:
        // TODO: Check for Go and TinyGo installations
        break;
    }

    // Verify and/or install the Modus SDK
    let installedVersion = await vi.getLatestInstalledVersion();
    const online = await isOnline();
    if (online) {
      let latestVersion: string | undefined;
      await withSpinner(chalk.dim("Checking to see if you have the latest Modus SDK version ..."), async () => {
        latestVersion = await vi.getLatestRuntimeVersion(prerelease);
        if (!latestVersion) {
          this.logError("Failed to fetch the latest Modus SDK version.");
          this.exit(1);
        }
      });

      let updateSDK = false;
      if (!installedVersion) {
        if (!(await this.confirmAction("You do not have the Modus SDK installed. Would you like to install it now? [y/n]"))) {
          clearLine();
          this.log(chalk.dim("Aborted."));
          this.exit(1);
        }
        updateSDK = true;
      } else if (latestVersion !== installedVersion) {
        if (await this.confirmAction("You have an outdated version of the Modus SDK. Would you like to update? [y/n]")) {
          updateSDK = true;
        } else {
          clearLine();
        }
      }
      if (updateSDK) {
        await SDKInstallCommand.run([latestVersion!]);
        installedVersion = latestVersion;
      }
    }

    if (!installedVersion) {
      this.logError("Could not find an installed Modus SDK.");
      this.exit(1);
    }

    const version = installedVersion;
    this.log(chalk.dim(`Using Modus SDK ${version}`));

    const templatesArchive = await vi.getLatestTemplatesArchivePath(version, sdk.toLowerCase());
    if (!templatesArchive) {
      this.logError(`Could not find any ${sdk} templates for SDK version ${version}`);
      this.exit(1);
    }

    // Install build tools if needed
    if (sdk == SDK.Go) {
      const ext = os.platform() === "win32" ? ".exe" : "";
      const buildTool = path.join(ModusHomeDir, "sdk", version, "modus-go-build" + ext);
      if (!(await fs.exists(buildTool))) {
        if (online) {
          const module = `github.com/${GitHubOwner}/${GitHubRepo}/sdk/go/tools/modus-go-build@${version}`;
          await withSpinner("Downloading the Modus Go build tool ...", async () => {
            await execFile("go", ["install", module], {
              cwd: ModusHomeDir,
              env: {
                ...process.env,
                GOBIN: path.join(ModusHomeDir, "sdk", version),
              },
            });
          });
        } else {
          this.logError("Could not find the Go build tool. Please try again when you are online.");
          this.exit(1);
        }
      }
    }

    // Create the app
    await withSpinner("Creating a new Modus app ...", async () => {
      if (!(await fs.exists(dir))) {
        await fs.mkdir(dir, { recursive: true });
      }
      await execFile("tar", ["-xf", templatesArchive, "-C", dir, "--strip-components=2", `templates/${template}`]);

      // Apply SDK-specific modifications
      const execOpts = { env: process.env, cwd: dir };
      switch (sdk) {
        case SDK.AssemblyScript:
          await execFile("npm", ["pkg", "set", `name=${name}`], execOpts);
          await execFile("npm", ["install"], execOpts);
          break;
        case SDK.Go:
          await execFile("go", ["mod", "download"], execOpts);
          break;
      }
    });

    this.log(chalk.bold(chalk.cyanBright("Successfully created a Modus app!")));
    this.log("To start, run the following command:");
    this.log("$ " + chalk.blueBright(`${dir == process.cwd() ? "" : "cd " + path.basename(dir)} && modus dev --build`));
    this.log("");
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }

  private async confirmAction(message: string): Promise<boolean> {
    this.log(message);
    const cont = ((await ask(chalk.dim(" -> "))) || "n").toLowerCase().trim();
    clearLine();
    return cont === "yes" || cont === "y";
  }
}
