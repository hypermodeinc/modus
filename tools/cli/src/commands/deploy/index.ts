import { Args, Command, Flags } from "@oclif/core";
import { isRunnable } from "../../util/index.js";
import chalk from "chalk";

export default class DeployCommand extends Command {
  static args = {
    path: Args.string()
  };
  static description = "Deploy a Modus app to Hypermode";
  static examples = [];
  static flags = {};

  async run(): Promise<void> {
    const { args, flags } = await this.parse(DeployCommand);

    if (!isHypCLIInstalled()) {
      this.logError("Hypermode CLI is not installed! Please run PLACEHOLDER to install and try again!");
      return;
    }

    
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

function isHypCLIInstalled(): boolean {
  return isRunnable("hyp");
}