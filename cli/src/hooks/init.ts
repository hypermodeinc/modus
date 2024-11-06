import { Hook } from "@oclif/core";

const hook: Hook.Init = async function (options) {
  if (["-v", "--version"].includes(process.argv[2])) {
    this.log(`v${options.config.version}`);
    process.exit(0);
  }
};

export default hook;
