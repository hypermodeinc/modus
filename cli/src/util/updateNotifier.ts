import chalk from "chalk";
import semver from "semver";
import * as vi from "./versioninfo.js";
import { isOnline } from "./index.js";

/**
 * Checks if there's a newer version of the CLI available and prints an update message if needed.
 * @param currentVersion Current CLI version
 */
export async function checkForUpdates(currentVersion: string): Promise<void> {
  if (!(await isOnline())) {
    return;
  }

  if (currentVersion.startsWith("v")) {
    currentVersion = currentVersion.slice(1);
  }

  const latestVersion = await vi.getLatestCliVersion(false);
  if (!latestVersion) return;

  const latestVersionNumber = latestVersion.startsWith("v") ? latestVersion.slice(1) : latestVersion;

  if (semver.gt(latestVersionNumber, currentVersion)) {
    console.log();
    console.log(chalk.yellow("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"));
    console.log(chalk.yellow(`Update available! ${chalk.dim(currentVersion)} → ${chalk.greenBright(latestVersionNumber)}`));
    console.log(chalk.yellow("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"));
    console.log();
  }
}
