import { Args, Command } from "@oclif/core";
import chalk from "chalk";

export default class SDKListCommand extends Command {
  static args = {};

  static description = "List installed SDK versions";

  static examples = ["modus list"];

  static flags = {};

  async run(): Promise<void> {
    const { args } = await this.parse(SDKListCommand);
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}
