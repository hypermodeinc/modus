import { Args, Command, Flags } from "@oclif/core";

export default class Run extends Command {
  static args = {
    watch: Args.string({
      aliases: ["w"],
      description: "Watch",
      required: false,
    }),
  };

  static description = "Run a Modus app locally";

  static examples = [`<%= config.bin %> <%= command.id %> run ./project-path --watch`];

  async run(): Promise<void> {
    const { args, flags } = await this.parse(Run);
  }
}
