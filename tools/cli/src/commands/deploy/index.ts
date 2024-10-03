import { Command, Flags } from "@oclif/core";

export default class DeployCommand extends Command {
  static args = {};
  static description = "Deploy a Modus app to Hypermode";
  static examples = [];
  static flags = {};

  async run(): Promise<void> {
    const { args, flags } = await this.parse(DeployCommand);
    
  }
}
