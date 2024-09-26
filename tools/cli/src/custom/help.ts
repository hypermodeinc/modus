import { Command, Help } from '@oclif/core'
import chalk from 'chalk';

import { VERSION } from "./globals.js";

export default class CustomHelp extends Help {
    // acts as a "router"
    // and based on the args it receives
    // calls one of showRootHelp, showTopicHelp,
    // the formatting for an individual command
    //   formatCommand(command: Command.Loadable): string {}

    // the formatting for a list of commands
    formatCommands(commands: Command.Loadable[]): string {
        let out = "";
        out += chalk.bold("Commands:") + "\n";
        for (const command of commands) {
            if (command.id === "autocomplete") continue;

            const name = chalk.bold.blueBright(command.id);
            const prePadding = " ".repeat(10 - command.id.length);
            const args = Object.keys(command.args).length > 0 ? Object.entries(command.args).map(v => {
                if (!v[1].hidden && v[1].required) {
                    if (v[1].description && v[1].description.indexOf("-|-") > 0) {
                        return v[1].description.split("-|-")[0];
                    }

                    return v[0];
                }

                return "";
            }).join(" ") : "";
            const postPadding = " ".repeat(Math.max(20 - args.length, 0));
            const description = command.description!;
            const aliases = command.aliases.length > 0 ? chalk.dim(" (" + command.aliases.join("/") + ")") : "";

            out += "  " + name + prePadding + chalk.dim(args) + postPadding + description + aliases + "\n";
        }

        out += "\n";
        out += "View the docs:                  " + chalk.blueBright("https://docs.hypermode.com/introduction") + "\n";
        out += "View the repo:                  " + chalk.blueBright("https://github.com/HypermodeAI/core") + "\n";

        out += "\n";
        out += "Made with 💖 by " + chalk.magentaBright("https://hypermode.ai");
        return out;
    }

    // displayed for the root help
    formatRoot(): string {
        let out = "";
        // I used to speak a little bit of finnish, but that says something like:
        // Hypermode is an awesome tool for ai development
        out += chalk.bold.blueBright("Hypermode") + " is kaunis tyokalu tekoalyn AI kehittamiseen! " + chalk.dim("(v" + VERSION + ")") + "\n\n";

        // Usage: hyp <command> [...flags] [...args]
        out += chalk.bold("Usage: hyp") + " " + chalk.dim("<command>") + " " + chalk.bold.blueBright("[...flags]") + " " + chalk.bold("[...args]");
        return out;
    }

    // the formatting for an individual topic
    //   formatTopic(topic: Interfaces.Topic): string {}

    // the default implementations of showRootHelp
    // showTopicHelp and showCommandHelp
    // will call various format methods that
    // provide the formatting for their corresponding
    // help sections;
    // these can be overwritten as well

    // the formatting responsible for the header
    // the formatting for a list of topics
    //   protected formatTopics(topics: Interfaces.Topic[]): string {}

    // This triggers when a user runs hyp <command> --help
    async showCommandHelp(command: Command.Loadable): Promise<void> {
        const name = command.id;
        const { examples } = command;

        const args = Object.keys(command.args).length > 0 ? Object.entries(command.args).map(v => {
            if (!v[1].hidden && v[1].required) {
                if (v[1].description && v[1].description.indexOf("-|-") > 0) {
                    return v[1].description.split("-|-")[0];
                }

                return v[0];
            }

            return "";
        }).join(" ") : "";
        const flags = Object.keys(command.flags).length > 0 ? " <flags>" : "";

        this.log(chalk.bold("Hypermode Help v0.0.0") + "\n");
        this.log(chalk.bold("Usage:") + "\n");

        this.log("  " + chalk.dim("$ hyp " + name + " " + args + flags));

        if (examples) {
            this.log();
            this.log(chalk.bold("Examples:") + "\n");
            for (const example of examples) this.log("  " + chalk.dim(example));
        }

        if (flags) {
            this.log();
            this.log(chalk.bold("Flags:") + "\n");
            for (const flag of Object.values(command.flags)) this.log("  " + chalk.dim(flag.name));
        }
    }

    // or showCommandHelp
    //   showHelp(args: string[]): void {}

    // display the root help of a CLI
    //   showRootHelp(): void {}

    // display help for a topic
    //   showTopicHelp(topic: Interfaces.Topic): void {}
}
