import { Command, Flags } from "@oclif/core"

export default class Upgrade extends Command {
  static args = {
    // person: Args.string({ description: "Person to say hello to", required: true }),
  }

  static description = "Upgrade a Hypermode component"

  static examples = [
    `<%= config.bin %> <%= command.id %> friend --from oclif
hello friend from oclif! (./src/commands/hello/index.ts)
`,
  ]

  static flags = {
    from: Flags.string({ char: "f", description: "Who is saying hello", required: true }),
  }

  async run(): Promise<void> {
    const { args, flags } = await this.parse(Upgrade)

    this.log(`hello ${args.person} from ${flags.from}! (./src/commands/hello/index.ts)`)
  }
}
