import { Args, Command, Flags } from "@oclif/core";
import { expandHomeDir } from "../../util/index.js";
import { Metadata } from "../../util/metadata.js";
import BuildCommand from "../build/index.js";
import path from "path";
import { copyFileSync, existsSync, readFileSync, watch as watchFolder } from "fs";
import chalk from "chalk";
import { spawn } from "child_process";
import os from "node:os";

export default class Run extends Command {
  static args = {
    path: Args.string({
      description: "./my-project-|-Directory to run",
      hidden: false,
      required: false,
    }),
  };

  static flags = {
    watch: Flags.boolean({
      description: "Watch project and rebuild continually",
      hidden: false,
      required: false,
    }),
  };

  static description = "Run a Modus app locally";

  static examples = [`<%= config.bin %> <%= command.id %> run ./project-path --watch`];

  async run(): Promise<void> {
    const { args, flags } = await this.parse(Run);
    const runtimePath = expandHomeDir("~/.hypermode/sdk/" + Metadata.runtime_version + "/runtime");

    const cwd = args.path ? path.join(process.cwd(), args.path) : process.cwd();
    const watch = flags.watch;

    if (!existsSync(path.join(cwd, "/package.json"))) {
      this.logError("Could not locate package.json! Please try again");
      process.exit(0);
    }

    let project_name: string;
    try {
      project_name = JSON.parse(readFileSync(path.join(cwd, "/package.json")).toString()).name;
    } catch {
      this.logError("Could not read package.json! Please try again");
      process.exit(0);
    }

    await BuildCommand.run(args.path ? [args.path] : []);
    const build_wasm = path.join(cwd, "/build/" + project_name + ".wasm");
    const deploy_wasm = expandHomeDir("~/.hypermode/" + project_name + ".wasm");
    copyFileSync(build_wasm, deploy_wasm);

    spawn(expandHomeDir("~/.hypermode/sdk/" + Metadata.runtime_version + "/runtime" + (os.platform() === "win32" ? ".exe" : "")), { stdio: "inherit" });

    if (watch) {
      const delay = 3000; // Max build frequency every 3000ms
      let lastModified = 0;
      let lastBuild = 0;
      let paused = true;

      // Its a bit jank. I'll refactor it in a bit
      setInterval(async () => {
        if (paused) return;
        if (lastBuild > lastModified) {
          paused = true;
          return;
        }
        if (Date.now() - lastModified > delay * 2) paused = true;
        lastBuild = Date.now();
        await BuildCommand.run(args.path ? [args.path] : []);
        copyFileSync(build_wasm, deploy_wasm);
      }, delay);

      watchFolder(path.join(cwd, "/assembly"), () => {
        lastModified = Date.now();
        paused = false;
      });
    }
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}
