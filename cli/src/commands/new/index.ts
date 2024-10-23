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
import * as inquirer from "@inquirer/prompts";

const MODUS_NEW_GO_NAME = "modus-go-app";
const MODUS_NEW_AS_NAME = "modus-as-app";

const MODUS_DEFAULT_TEMPLATE_NAME = "default";

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
    // template: Flags.string({
    //   char: "t",
    //   description: "Template to use",
    // }),
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

    let sdk: SDK;
    if (flags.sdk) {
      sdk = parseSDK(flags.sdk);
    } else {
      const sdkInput = await inquirer.select({
        message: "Select a SDK",
        default: SDK.Go,
        choices: [
          {
            value: SDK.Go,
          },
          {
            value: SDK.AssemblyScript,
          },
        ],
      });

      sdk = parseSDK(sdkInput);
    }

    const defaultName = sdk === SDK.Go ? MODUS_NEW_GO_NAME : MODUS_NEW_AS_NAME;

    const name =
      flags.name ||
      (await inquirer.input({
        message: "Pick a name for your app:",
        default: defaultName,
      }));

    const dir = flags.dir || "." + path.sep + name;

    if (!flags.force) {
      const confirm = await inquirer.confirm({ message: "Continue?", default: true });
      if (!confirm) {
        this.log(chalk.dim("Aborted"));
        this.exit(1);
      }
    }

    this.log();
    await this.createApp(name, dir, sdk, MODUS_DEFAULT_TEMPLATE_NAME, flags.force, flags.prerelease);
  }

  private async createApp(name: string, dir: string, sdk: SDK, template: string, force: boolean, prerelease: boolean) {
    if (!force && (await fs.exists(dir))) {
      if (!(await this.confirmAction("Attempting to overwrite a folder that already exists.\nAre you sure you want to continue? [Y/n]"))) {
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
        const tinyGoVersion = await getTinyGoVersion();

        const foundGo = !!goVersion;
        const foundTinyGo = !!tinyGoVersion;
        const okGo = foundGo && semver.gte(goVersion, MinGoVersion);
        const okTinyGo = foundTinyGo && semver.gte(tinyGoVersion, MinTinyGoVersion);

        if (okGo && okTinyGo) {
          break;
        }

        this.log(`The Modus Go SDK requires both of the following:`);

        if (okGo) {
          this.log(chalk.dim(`• Go v${MinGoVersion} or newer `) + chalk.green(`(found v${goVersion})`));
        } else if (foundGo) {
          this.log(chalk.dim(`• Go v${MinGoVersion} or newer `) + chalk.red(`(found v${goVersion})`));
        } else {
          this.log(chalk.dim(`• Go v${MinGoVersion} or newer `) + chalk.red("(not found)"));
        }

        if (okTinyGo) {
          this.log(chalk.dim(`• TinyGo v${MinTinyGoVersion} or newer `) + chalk.green(`(found v${tinyGoVersion})`));
        } else if (foundTinyGo) {
          this.log(chalk.dim(`• TinyGo v${MinTinyGoVersion} or newer `) + chalk.red(`(found v${tinyGoVersion})`));
        } else {
          this.log(chalk.dim(`• TinyGo v${MinTinyGoVersion} or newer `) + chalk.red("(not found)"));
        }

        this.log();

        if (!okGo) {
          this.log(chalk.yellow(`Please install Go ${MinGoVersion} or newer from https://go.dev/dl`));
          this.log();
        }

        if (!okTinyGo) {
          this.log(chalk.yellow(`Please install TinyGo ${MinTinyGoVersion} or newer from https://tinygo.org/getting-started/install`));
          if (os.platform() === "win32") {
            this.log(`Note that you will need to install the binaryen components for wasi support.`);
          }
          this.log();
        }

        this.log(`Make sure to add Go and TinyGo to your PATH. You can check for yourself by running:`);
        this.log(`${chalk.dim("$")} go version`);
        this.log(`${chalk.dim("$")} tinygo version`);
        this.log();

        this.exit(1);
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
        if (!(await this.confirmAction(`You do not have the ${sdkText} installed. Would you like to install it now? [Y/n]`))) {
          this.log(chalk.dim("Aborted."));
          this.exit(1);
        }
        updateSDK = true;
      } else if (latestVersion !== installedSdkVersion) {
        if (await this.confirmAction(`You have ${installedSdkVersion} of the ${sdkText}. The latest is ${latestVersion}. Would you like to update? [Y/n]`)) {
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

  private async confirmAction(message: string, defaultToContinue = true): Promise<boolean> {
    this.log(message);
    const input = await ask(chalk.dim(" -> "));

    if (input === "") {
      clearLine(2);
      return defaultToContinue;
    }

    const shouldContinue = (input || "n").toLowerCase().trim();
    clearLine(2);
    return shouldContinue === "yes" || shouldContinue === "y";
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

function toValidAppName(input: string): string {
  // Remove any characters that aren't alphanumeric, spaces, or a few other valid characters.
  // Replace spaces with hyphens.
  return input
    .trim() // Remove leading/trailing spaces
    .toLowerCase() // Convert to lowercase for consistency
    .replace(/[^a-z0-9\s-]/g, "") // Remove invalid characters
    .replace(/\s+/g, "-") // Replace spaces (or multiple spaces) with a single hyphen
    .replace(/-+/g, "-") // Replace multiple consecutive hyphens with a single hyphen
    .replace(/^-|-$/g, ""); // Remove leading or trailing hyphens
}
