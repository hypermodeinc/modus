import {Args, Command, Flags} from "@oclif/core"

export default class HelpCommand extends Command {
  static args = {
    person: Args.string({description: "<command>-|-Print help for command", required: true}),
  }

  static description = "Print help for certain command"

  static examples = [
    `<%= config.bin %> <%= command.id %> friend --from oclif
hello friend from oclif! (./src/commands/hello/index.ts)
`,
  ]

  static flags = {
    from: Flags.string({char: "f", description: "Who is saying hello", required: true}),
  }

  async run(): Promise<void> {
    const {args, flags} = await this.parse(HelpCommand)

    this.log(`hello ${args.person} from ${flags.from}! (./src/commands/hello/index.ts)`)
  }
}
