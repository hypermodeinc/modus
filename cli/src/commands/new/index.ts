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
import semver from "semver";
import os from "node:os";
import path from "node:path";

import * as fs from "../../util/fs.js";
import * as vi from "../../util/versioninfo.js";
import { execFile } from "../../util/cp.js";
import { isOnline } from "../../util/index.js";
import { GitHubOwner, GitHubRepo, MinGoVersion, MinNodeVersion, MinTinyGoVersion, ModusHomeDir, SDK, parseSDK } from "../../custom/globals.js";
import { ask, clearLine, withSpinner } from "../../util/index.js";
import SDKInstallCommand from "../sdk/install/index.js";
import { getHeader } from "../../custom/header.js";

export default class NewCommand extends Command {
  static description = "Create a new Modus app";

  static examples = ["modus new", "modus new --name my-app", "modus new --name my-app --sdk go --dir ./my-app --force"];

  static flags = {
    nologo: Flags.boolean({
      aliases: ["no-logo"],
      hidden: true,
    }),
    name: Flags.string({
      char: "n",
      description: "App name",
    }),
    dir: Flags.directory({
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

    if (!flags.nologo) {
      this.log(getHeader(this.config.version));
    }

    this.log(chalk.hex("#A585FF")(NewCommand.description) + "\n");
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
    if (!flags.force && !(await this.confirmAction("Continue? [y/n]"))) {
      this.log(chalk.dim("Aborted."));
      this.exit(1);
    }
    this.log();
    await this.createApp(name, dir, sdk, template, flags.force, flags.prerelease);
  }

  private async promptAppName(): Promise<string> {
    this.log("App Name?");
    const name = ((await ask(chalk.dim(" -> "))) || "").trim();
    clearLine(2);
    this.log("App Name: " + chalk.dim(name.length ? name : "Not Provided"));
    return name;
  }

  private async promptInstallPath(defaultValue: string): Promise<string> {
    this.log("Install Directory? " + chalk.dim(`(${defaultValue})`));
    const dir = ((await ask(chalk.dim(" -> "))) || defaultValue).trim();
    clearLine(2);
    this.log("Directory: " + chalk.dim(dir));
    return path.resolve(dir);
  }

  private async promptSdkSelection(): Promise<string> {
    this.log("Select an SDK");
    for (const [index, sdk] of Object.values(SDK).entries()) {
      this.log(chalk.dim(` ${index + 1}. ${sdk}`));
    }

    const selectedIndex = Number.parseInt(((await ask(chalk.dim(" -> "))) || "1").trim(), 10) - 1;
    const sdk = Object.values(SDK)[selectedIndex];
    clearLine(Object.values(SDK).length + 2);
    if (!sdk) this.exit(1);
    this.log("SDK: " + chalk.dim(sdk));
    return sdk;
  }

  private async promptTemplate(defaultValue: string): Promise<string> {
    this.log("Template? " + chalk.dim(`(${defaultValue})`));
    const template = ((await ask(chalk.dim(" -> "))) || defaultValue).trim();
    clearLine(2);
    this.log("Template: " + chalk.dim(template));
    return template;
  }

  private async createApp(name: string, dir: string, sdk: SDK, template: string, force: boolean, prerelease: boolean) {
    if (!force && (await fs.exists(dir))) {
      if (!(await this.confirmAction("Attempting to overwrite a folder that already exists.\nAre you sure you want to continue? [y/n]"))) {
        clearLine();
        return;
      } else {
        clearLine(2);
      }
    }

    // Validate SDK-specific prerequisites
    const sdkText = `Modus ${sdk} SDK`;
    switch (sdk) {
      case SDK.AssemblyScript:
        if (semver.lt(process.versions.node, MinNodeVersion)) {
          this.logError(`The ${sdkText} requires Node.js version ${MinNodeVersion} or newer.`);
          this.logError(`You have Node.js version ${process.versions.node}.`);
          this.logError(`Please upgrade Node.js and try again.`);
          this.exit(1);
        }
        break;
      case SDK.Go:
        const goVersion = await getGoVersion();
        if (!goVersion) {
          this.logError(`Go is not installed. Please install Go ${MinGoVersion} or newer and try again.`);
          this.logError(`You will also need to install TinyGo ${MinTinyGoVersion} or newer.`);
          this.exit(1);
        }
        if (semver.lt(goVersion, MinGoVersion)) {
          this.logError(`The ${sdkText} requires Go version ${MinGoVersion} or newer.`);
          this.logError(`You have Go version ${goVersion}.`);
          this.logError(`Please upgrade Go and try again.`);
          this.exit(1);
        }

        const tinyGoVersion = await getTinyGoVersion();
        if (!tinyGoVersion) {
          this.logError(`TinyGo is not installed. Please install Go ${MinTinyGoVersion} or newer and try again.`);
          this.exit(1);
        }
        if (semver.lt(tinyGoVersion, MinTinyGoVersion)) {
          this.logError(`The Modus Go SDK requires TinyGo version ${MinTinyGoVersion} or newer.`);
          this.logError(`You have TinyGo version ${tinyGoVersion}.`);
          this.logError(`Please upgrade TinyGo and try again.`);
          this.exit(1);
        }

        break;
    }

    // Verify and/or install the Modus SDK
    let installedSdkVersion = await vi.getLatestInstalledSdkVersion(sdk, prerelease);
    if (await isOnline()) {
      let latestVersion: string | undefined;
      await withSpinner(chalk.dim(`Checking to see if you have the latest version of the ${sdkText}.`), async () => {
        latestVersion = await vi.getLatestSdkVersion(sdk, prerelease);
        if (!latestVersion) {
          this.logError(`Failed to retrieve the latest ${sdkText} version.`);
          this.exit(1);
        }
      });

      let updateSDK = false;
      if (!installedSdkVersion) {
        if (!(await this.confirmAction(`You do not have the ${sdkText} installed. Would you like to install it now? [y/n]`))) {
          this.log(chalk.dim("Aborted."));
          this.exit(1);
        }
        updateSDK = true;
      } else if (latestVersion !== installedSdkVersion) {
        if (await this.confirmAction(`You have ${installedSdkVersion} of the ${sdkText}. The latest is ${latestVersion}. Would you like to update? [y/n]`)) {
          updateSDK = true;
        }
      }
      if (updateSDK) {
        await SDKInstallCommand.run([sdk, latestVersion!, "--no-logo"]);
        installedSdkVersion = latestVersion;
      }
    }

    if (!installedSdkVersion) {
      this.logError(`Could not find an installed ${sdkText}.`);
      this.exit(1);
    }

    const sdkVersion = installedSdkVersion;
    const sdkPath = vi.getSdkPath(sdk, sdkVersion);

    const templatesArchive = path.join(sdkPath, "templates.tar.gz");
    if (!(await fs.exists(templatesArchive))) {
      this.logError(`Could not find any templates for ${sdkText} ${sdkVersion}`);
      this.exit(1);
    }

    // Install build tools if needed
    if (sdk == SDK.Go) {
      const ext = os.platform() === "win32" ? ".exe" : "";
      const buildTool = path.join(sdkPath, "modus-go-build" + ext);
      if (!(await fs.exists(buildTool))) {
        if (await isOnline()) {
          const module = `github.com/${GitHubOwner}/${GitHubRepo}/sdk/go/tools/modus-go-build@${sdkVersion}`;
          await withSpinner("Downloading the Modus Go build tool.", async () => {
            await execFile("go", ["install", module], {
              cwd: ModusHomeDir,
              env: {
                ...process.env,
                GOBIN: sdkPath,
              },
            });
          });
        } else {
          this.logError("Could not find the Modus Go build tool. Please try again when you are online.");
          this.exit(1);
        }
      }
    }

    // Create the app
    this.log(chalk.dim(`Using ${sdkText} ${sdkVersion}`));
    await withSpinner(`Creating a new Modus ${sdk} app.`, async () => {
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

    this.log();
    this.log(chalk.bold.cyanBright(`Successfully created a Modus ${sdk} app!`));
    this.log();

    this.log("To start, run the following command:");
    this.log("$ " + chalk.blueBright(`${dir == process.cwd() ? "" : "cd " + path.basename(dir)} && modus dev`));
    this.log();
  }

  private logError(message: string) {
    this.log(chalk.red(" ERROR ") + chalk.dim(": " + message));
  }

  private async confirmAction(message: string): Promise<boolean> {
    this.log(message);
    const cont = ((await ask(chalk.dim(" -> "))) || "n").toLowerCase().trim();
    clearLine(2);
    return cont === "yes" || cont === "y";
  }
}

async function getGoVersion(): Promise<string | undefined> {
  try {
    const result = await execFile("go", ["version"], {
      cwd: ModusHomeDir,
      env: process.env,
    });
    const parts = result.stdout.split(" ");
    const str = parts.length > 2 ? parts[2] : undefined;
    if (str?.startsWith("go")) {
      return str.slice(2);
    }
  } catch {}
}

async function getTinyGoVersion(): Promise<string | undefined> {
  try {
    const result = await execFile("tinygo", ["version"], {
      cwd: ModusHomeDir,
      env: process.env,
    });
    const parts = result.stdout.split(" ");
    return parts.length > 2 ? parts[2] : undefined;
  } catch {}
}
