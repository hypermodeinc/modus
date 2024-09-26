import { Command, Flags} from "@oclif/core"

export default class DeployCommand extends Command {
  static args = {}

  static description = "Deploy a Hypermode app to GitHub"

  static examples = [
    `<%= config.bin %> <%= command.id %> friend --from oclif
hello friend from oclif! (./src/commands/hello/index.ts)
`,
  ]

  static flags = {
    from: Flags.string({char: "f", description: "Who is saying hello", required: true}),
  }

  async run(): Promise<void> {
    const {args, flags} = await this.parse(DeployCommand)

    this.log(`hello ${args.person} from ${flags.from}! (./src/commands/hello/index.ts)`)
  }
}
