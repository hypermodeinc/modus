/* eslint-disable n/no-process-exit */
/* eslint-disable unicorn/no-process-exit */
import { Command, Flags } from "@oclif/core";
import sleep from "atomic-sleep";
import chalk from "chalk";
import { createInterface } from "node:readline";
import ora from "ora";

import { SDK, VERSION } from "../../custom/globals.js";
import {
  ask,
  clearLine,
  cloneRepo,
  getAvailablePackageManagers,
  installDeps,
  isGitInstalled,
  isGoInstalled,
  isTinyGoInstalled,
} from "../../util/index.js";
import Runtime from "../../util/runtime.js";
const PKGMGRS = getAvailablePackageManagers();

export default class NewCommand extends Command {
  static description = "Create a new Hypermode project";

  static examples = [`<%= config.bin %> <%= command.id %> --sdk go`];

  static flags = {
    name: Flags.string({ description: "Project name" }),
    pkgmgr: Flags.string({
      description: "Package manager to use",
    }),
    sdk: Flags.string({
      description: "SDK to use",
    })
  };

  async run(): Promise<void> {
    const { flags } = await this.parse(NewCommand);

    const rl = createInterface({
      input: process.stdin,
      output: process.stdout,
    });

    this.log(chalk.bold(`Hypermode new v${VERSION}\n`));

    if (PKGMGRS.length === 0) {
      this.logError("Could not find any suitable package manager. Please install NPM, Yarn, PNPM, or Bun!");
      return;
    }

    const name = flags.name || await this.promptProjectName(rl);
    const sdk = flags.sdk ? Object.values(SDK)[Object.keys(SDK).map(v => v.toLowerCase()).indexOf(flags.sdk?.trim().toLowerCase())] : await this.promptSdkSelection(rl); // Use the enum
    // Don't use during init, perhaps use a flag or not
    const pkgMgr = flags.pkgmgr || await this.promptPackageManager(sdk, rl);

    if (!(flags.name && flags.sdk && (flags.sdk.trim().toLowerCase() === "assemblyscript" && flags.pkgmgr)) && !(await this.confirmAction(rl, `[4/4] Continue? [y/n]`))) process.exit(0);
    await this.installProject(name, sdk, pkgMgr);
  }

  private async checkGoInstallation() {
    if (!isGoInstalled()) {
      this.logError("Could not find valid Go installation! Please install Go and try again.");
      process.exit(0);
    }

    if (!isTinyGoInstalled()) {
      this.logError("Could not find valid TinyGo installation! Please install TinyGo and try again.");
      process.exit(0);
    }
  }

  private async confirmAction(rl: ReturnType<typeof createInterface>, message: string): Promise<boolean> {
    this.log(message);
    const cont = (await ask(chalk.dim(" -> "), rl)).toLowerCase().trim() === "y";
    clearLine();
    return cont;
  }

  private async downloadAndInstallRuntime() {
    const latestRelease = await Runtime.getLatestRelease();
    const runtimeDlSpinner = ora({ color: "white", indent: 2, text: `Downloading Runtime ${chalk.dim(`(${latestRelease})`)}` }).start();
    sleep(2000); // Runtime is private on gh, soooo can't really do anything w/o ssh
    runtimeDlSpinner.stop();

    const runtimeInstSpinner = ora({ color: "white", indent: 2, text: `Installing Runtime ${chalk.dim(`(${latestRelease})`)}` }).start();
    sleep(2000);
    runtimeInstSpinner.stop();
    Runtime.install(latestRelease!);
    this.log(`- Installed Runtime ${chalk.dim(`(${latestRelease})`)}`);

    Runtime.install("0.12.0");
  }

  private async installProject(name: string, sdk: string, pkgMgr: string) {
    this.log("[4/4] Installing");

    if (!isGitInstalled()) {
      this.logError("Could not find valid Git installation! Please download Git or ensure it is in your PATH!");
      return;
    }

    const gitSpinner = ora({ color: "white", indent: 2, text: "Cloning repo" }).start();
    const cloned = cloneRepo("https://github.com/hypermodeAI/template-project", name);

    if (!cloned) {
      gitSpinner.stop();
      this.logError("Failed to clone the git repository. Please check your internet and try again.");
      return;
    }

    gitSpinner.stop();
    this.log("- Cloned Repo");

    const depsSpinner = ora({ color: "white", indent: 2, text: "Installing dependencies" }).start();
    const depsInstalled = installDeps(pkgMgr, name);

    if (!depsInstalled) {
      depsSpinner.stop();
      this.logError("Failed to install dependencies! Please try again.");
      return;
    }

    depsSpinner.stop();
    this.log("- Installed Dependencies");

    await this.downloadAndInstallRuntime();

    this.log("\nSuccessfully installed Hypermode SDK!");
    this.log("To start, run the following command:");
    this.log(chalk.dim(`$ cd ${name} && hyp run`));
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }

  private async promptPackageManager(sdk: string, rl: ReturnType<typeof createInterface>): Promise<string> {
    if (sdk === SDK.Go) {
      await this.checkGoInstallation();
      return "go";
    }

    this.log("[3/4] Select package manager");
    for (const [index, mgr] of PKGMGRS.entries()) this.log(chalk.dim(` ${index + 1}. ${mgr}`));
    const selectedIndex = Number.parseInt((await ask(chalk.dim(" -> "), rl)).trim(), 10) - 1;
    const pkgMgr = PKGMGRS[selectedIndex] || "NPM"; clearLine(); clearLine();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    for (const _ of PKGMGRS) clearLine();
    this.log("[3/4] Package Manager: " + chalk.dim(pkgMgr));
    return pkgMgr;
  }

  private async promptProjectName(rl: ReturnType<typeof createInterface>): Promise<string> {
    this.log("[1/4] Project Name:");
    const name = (await ask(chalk.dim(" -> "), rl)).trim(); clearLine(); clearLine();
    this.log("[1/4] Name: " + chalk.dim(name));
    return name;
  }

  private async promptSdkSelection(rl: ReturnType<typeof createInterface>): Promise<string> {
    this.log("[2/4] Select an SDK");
    for (const [index, sdk] of Object.values(SDK).entries()) {
      this.log(chalk.dim(` ${index + 1}. ${sdk}`));
    }

    const selectedIndex = Number.parseInt((await ask(chalk.dim(" -> "), rl)).trim(), 10) - 1;
    const sdk = Object.values(SDK)[selectedIndex]; clearLine(); clearLine();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    for (const _ of Object.values(SDK)) clearLine();
    if (!sdk) process.exit(1);
    this.log("[2/4] SDK: " + chalk.dim(sdk));
    return sdk;
  }
}
