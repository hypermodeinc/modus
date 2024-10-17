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

import { createInterface } from "node:readline";
import path from "node:path";

import * as fs from "../../util/fs.js";
import { execFile } from "../../util/cp.js";
import { SDK } from "../../custom/globals.js";
import { ask, clearLine, expandHomeDir } from "../../util/index.js";
import * as vi from "../../util/versioninfo.js";

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
      description: "Initialize without prompt",
    }),
  };

  async run(): Promise<void> {
    const { flags } = await this.parse(NewCommand);

    const rl = createInterface({
      input: process.stdin,
      output: process.stdout,
    });

    this.log(chalk.bold(`Modus CLI v${this.config.version}\n${flags.force ? chalk.dim("WARN: Running in forced mode! Proceed with caution.") : ""}`));

    const name = flags.name || (await this.promptAppName(rl));
    if (!name) {
      this.logError("An app name is required.");
      this.exit(1);
    }

    const dir = flags.dir ? path.join(process.cwd(), flags.dir) : await this.promptInstallPath(rl, "." + path.sep + name);
    if (!dir) {
      this.logError("An install directory is required.");
      this.exit(1);
    }

    let sdk = flags.sdk
      ? Object.values(SDK)[
          Object.keys(SDK)
            .map((v) => v.toLowerCase())
            .indexOf(flags.sdk?.trim().toLowerCase())
        ]
      : await this.promptSdkSelection(rl); // Use the enum
    sdk = sdk.toLowerCase();

    const template = flags.template || (await this.promptTemplate(rl, "default"));
    if (!template) {
      this.logError("A template is required.");
      this.exit(1);
    }

    if (!flags.force && !(await this.confirmAction(rl, "[5/5] Continue? [y/n]"))) clearLine(), clearLine(), process.exit(0);

    await this.createApp(name, dir, sdk, template, flags.force, rl);
    // this.exit(0);
  }

  private async promptAppName(rl: ReturnType<typeof createInterface>): Promise<string> {
    this.log("[1/5] App Name:");
    const name = ((await ask(chalk.dim(" -> "), rl)) || "").trim();
    clearLine();
    clearLine();
    this.log("[1/5] Name: " + chalk.dim(name.length ? name : "Not Provided"));
    return name;
  }

  private async promptInstallPath(rl: ReturnType<typeof createInterface>, defaultValue: string): Promise<string> {
    this.log("[2/5] Install Dir: " + chalk.dim(`(${defaultValue})`));
    const dir = ((await ask(chalk.dim(" -> "), rl)) || defaultValue).trim();
    clearLine();
    clearLine();
    this.log("[2/5] Directory: " + chalk.dim(dir));
    return path.resolve(dir);
  }

  private async promptSdkSelection(rl: ReturnType<typeof createInterface>): Promise<string> {
    this.log("[3/5] Select an SDK");
    for (const [index, sdk] of Object.values(SDK).entries()) {
      this.log(chalk.dim(` ${index + 1}. ${sdk}`));
    }

    const selectedIndex = Number.parseInt(((await ask(chalk.dim(" -> "), rl)) || "1").trim(), 10) - 1;
    const sdk = Object.values(SDK)[selectedIndex];
    clearLine();
    clearLine();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    for (const _ of Object.values(SDK)) clearLine();
    if (!sdk) process.exit(1);
    this.log("[3/5] SDK: " + chalk.dim(sdk));
    return sdk;
  }

  private async promptTemplate(rl: ReturnType<typeof createInterface>, defaultValue: string): Promise<string> {
    this.log("[4/5] Template: " + chalk.dim(`(${defaultValue})`));
    const template = ((await ask(chalk.dim(" -> "), rl)) || defaultValue).trim();
    clearLine();
    clearLine();
    this.log("[4/5] Template: " + chalk.dim(template));
    return template;
  }

  private async createApp(name: string, dir: string, sdk: string, template: string, force: boolean, rl: ReturnType<typeof createInterface>) {
    if (!force && (await fs.exists(dir))) {
      if (!(await this.confirmAction(rl, "Attempting to overwrite a folder that already exists.\nAre you sure you want to continue? [y/n]"))) {
        clearLine();
        process.exit(0);
      } else {
        clearLine(), clearLine();
      }
    }

    clearLine();
    this.log("[3/4] Installing");

    // TODO: validate prerequisites for the SDK

    if (!(await fs.exists(dir))) {
      await fs.mkdir(dir, { recursive: true });
    }

    const version = vi.latestInstalledVersion();
    if (!version) {
      // TODO: install the latest version
      this.logError("Could not find the latest installed SDK version.");
      process.exit(1);
    }

    const templatesArchive = vi.getLatestTemplatesArchivePath(version, sdk);
    if (!templatesArchive) {
      this.logError(`Could not find any templates for SDK version ${version}`);
      process.exit(1);
    }

    await execFile("tar", ["-xf", templatesArchive, "-C", dir, "--strip-components=2", `templates/${template}`]);

    // if (!existsSync(templatePath)) {
    //   this.logError("Could not find the template for the latest installed SDK version.");
    //   process.exit(1);
    // }
    // cpSync(templatePath, dir, { recursive: true });

    // const depsSpinner = ora({
    //   color: "white",
    //   indent: 2,
    //   text: "Installing dependencies",
    // }).start();

    // if (sdk === "AssemblyScript") {
    //   if (isRunnable(NPM_CMD)) execSync(quote([NPM_CMD, "install"]), { cwd: dir, stdio: "ignore" });
    // } else if (sdk === "Go (Beta)") {
    //   const sh = execSync("go install", { cwd: dir, stdio: "ignore" });
    //   if (!sh) {
    //     this.logError("Failed to install dependencies via go install! Please try again");
    //     process.exit(0);
    //   }
    // }

    // depsSpinner.stop();
    // this.log("- Installed Dependencies");

    // await this.installRuntime();

    this.log("\nSuccessfully created a Modus app!");
    this.log("To start, run the following command:");
    this.log("$ " + chalk.dim(`${dir == process.cwd() ? "" : "cd " + path.basename(dir)} && modus dev --build`));
  }

  // private async installRuntime() {
  //   const latest_runtime = await getLatestRuntimeVersion(true);

  //   if (!latest_runtime) {
  //     this.logError("Could not find latest runtime via GitHub API. Please try again with internet access!");
  //     process.exit(0);
  //   }

  //   if (!Metadata.runtimes.includes(latest_runtime)) {
  //     const runtimeDlSpinner = ora({
  //       color: "white",
  //       indent: 2,
  //       text: `Downloading Runtime ${chalk.dim(`(${latest_runtime})`)}`,
  //     }).start();
  //     runtimeDlSpinner.stop();

  //     const runtimeInstSpinner = ora({
  //       color: "white",
  //       indent: 2,
  //       text: `Installing Runtime ${chalk.dim(`(${latest_runtime})`)}`,
  //     }).start();
  //     runtimeInstSpinner.stop();

  //     SDKInstallCommand.run([latest_runtime, "--silent"]);
  //   }
  //   this.log(`- Installed Runtime ${chalk.dim(`(${latest_runtime})`)}`);
  // }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }

  private async confirmAction(rl: ReturnType<typeof createInterface>, message: string): Promise<boolean> {
    this.log(message);
    const cont = ((await ask(chalk.dim(" -> "), rl)) || "n").toLowerCase().trim();
    clearLine();
    return cont === "yes" || cont === "y";
  }
}
