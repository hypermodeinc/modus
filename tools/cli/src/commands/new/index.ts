import { Command, Flags } from "@oclif/core";
import sleep from "atomic-sleep";
import chalk from "chalk";
import { createInterface } from "node:readline";
import ora from "ora";

import { CLI_VERSION, SDK } from "../../custom/globals.js";
import { ask, clearLine, cloneRepo, expandHomeDir, getAvailablePackageManagers, isGitInstalled } from "../../util/index.js";
import path from "node:path";
import { Metadata } from "../../util/metadata.js";
import { copyFileSync, existsSync, mkdirSync } from "node:fs";

const PKGMGRS = getAvailablePackageManagers();

export default class NewCommand extends Command {
  static description = "Create a new Modus project";

  static examples = ["modus new", "modus new --name Project01", "modus new --name Project01 --sdk go --dir ./my-project --force"];

  static flags = {
    name: Flags.string({ description: "Project name" }),
    dir: Flags.string({
      description: "Directory to install to",
      aliases: ["d"],
    }),
    sdk: Flags.string({ description: "SDK to use" }),
    force: Flags.boolean({
      description: "Initialize without prompt",
      aliases: ["f"],
    }),
  };

  async run(): Promise<void> {
    const { flags } = await this.parse(NewCommand);

    const rl = createInterface({
      input: process.stdin,
      output: process.stdout,
    });

    this.log(chalk.bold(`Modus new v${CLI_VERSION}\n${flags.force ? chalk.dim("WARN: Running in forced mode! Proceed with caution.") : ""}`));

    if (PKGMGRS.length === 0) {
      this.logError("Could not find any suitable package manager. Please install NPM, Yarn, PNPM, or Bun!");
      return;
    }

    const name = flags.name || (await this.promptProjectName(rl));
    const dir = flags.dir ? path.join(process.cwd(), flags.dir) : await this.promptInstallPath(rl);
    const sdk = flags.sdk
      ? Object.values(SDK)[
          Object.keys(SDK)
            .map((v) => v.toLowerCase())
            .indexOf(flags.sdk?.trim().toLowerCase())
        ]
      : await this.promptSdkSelection(rl); // Use the enum

    if (!flags.force && !(await this.confirmAction(rl, "[3/4] Continue? [y/n]"))) clearLine(), clearLine(), process.exit(0);

    this.installProject(name, dir, sdk, flags.force, rl);
  }

  private async promptProjectName(rl: ReturnType<typeof createInterface>): Promise<string> {
    this.log("[1/4] Project Name:");
    const name = ((await ask(chalk.dim(" -> "), rl)) || "").trim();
    clearLine();
    clearLine();
    this.log("[1/4] Name: " + chalk.dim(name.length ? name : "Not Provided"));
    return name;
  }

  private async promptInstallPath(rl: ReturnType<typeof createInterface>): Promise<string> {
    this.log("[2/4] Install Dir:");
    const dir = ((await ask(chalk.dim(" -> "), rl)) || "./").trim();
    clearLine();
    clearLine();
    this.log("[2/4] Directory: " + chalk.dim(dir));
    return path.join(process.cwd(), dir);
  }

  private async promptSdkSelection(rl: ReturnType<typeof createInterface>): Promise<string> {
    this.log("[2/4] Select an SDK");
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
    this.log("[2/4] SDK: " + chalk.dim(sdk));
    return sdk;
  }

  private async installProject(name: string, dir: string, sdk: string, force: boolean, rl: ReturnType<typeof createInterface>) {
    if (!force && !existsSync(dir)) {
      if (!(await this.confirmAction(rl, "Attempting to overwrite a folder that already exists.\nAre you sure you want to continue? [y/n]"))) clearLine(), process.exit(0);
      else clearLine(), clearLine(), clearLine();
    }

    this.log("[3/4] Installing");

    if (!isGitInstalled()) {
      this.logError("Could not find valid Git installation! Please download Git or ensure it is in your PATH!");
      return;
    }

    const gitSpinner = ora({
      color: "white",
      indent: 2,
      text: "Cloning template repository",
    }).start();
    const clone = cloneRepo("https://github.com/HypermodeAI/template-project", dir);

    if (!clone) {
      gitSpinner.stop();
      this.logError("Failed to clone the git repository. Please check your internet and try again.");
      return;
    }

    gitSpinner.stop();
    this.log("- Cloned template repository");

    await this.installRuntime();

    this.log("\nSuccessfully installed Modus SDK!");
    this.log("To start, run the following command:");
    this.log(chalk.dim(`$ ${dir == process.cwd() ? "" : "cd " + path.basename(dir)} && modus run`));
  }

  // TODO: install deps

  private async installRuntime() {
    const latest_runtime = await Metadata.getLatestRuntime();

    if (!latest_runtime) {
      this.logError("Could not find latest runtime via GitHub API. Please try again with internet access!");
      process.exit(0);
    }

    if (!Metadata.runtimes.includes(latest_runtime)) {
      const runtimeDlSpinner = ora({
        color: "white",
        indent: 2,
        text: `Downloading Runtime ${chalk.dim(`(${latest_runtime})`)}`,
      }).start();
      sleep(2000); // Runtime is private on gh, soooo can't really do anything w/o ssh
      runtimeDlSpinner.stop();

      const runtimeInstSpinner = ora({
        color: "white",
        indent: 2,
        text: `Installing Runtime ${chalk.dim(`(${latest_runtime})`)}`,
      }).start();
      sleep(2000);
      runtimeInstSpinner.stop();

      const outDir = expandHomeDir("~/.hypermode/sdk/" + latest_runtime.replace("v", ""));
      mkdirSync(outDir, { recursive: true });
      copyFileSync(path.join(process.cwd(), "./runtime-bin/" + latest_runtime.replace("v", "")), outDir + "/runtime");
    }
    this.log(`  - Installed Runtime ${chalk.dim(`(${latest_runtime})`)}`);
  }

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
