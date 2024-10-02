import { Args, Command } from "@oclif/core";
import chalk from "chalk";
import { readdirSync } from "node:fs";
import { expandHomeDir } from "../../../util/index.js";

export default class SDKListCommand extends Command {
  static args = {};
  static description = "List installed SDK versions";
  static examples = ["modus sdk list"];
  static flags = {};

  async run(): Promise<void> {
    const { args } = await this.parse(SDKListCommand);

    let versions: string[] = [];

    try {
      versions = readdirSync(expandHomeDir("~/.hypermode/sdk")).reverse();
    } catch {
      versions = [];
    }

    if (!versions.length) {
      this.log("No runtimes installed!");
      return;
    }

    for (const version of versions) {
      this.log(version);
    }
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}
