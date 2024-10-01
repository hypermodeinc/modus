import { execSync, spawnSync } from "node:child_process";
import { cpSync, existsSync, mkdirSync, rmSync } from "node:fs";
import path from "node:path";
import { Interface } from "node:readline";
import { CLI_VERSION } from "../custom/globals.js";

export async function ensureDir(dir: string): Promise<void> {
  if (!existsSync(dir)) mkdirSync(dir, { recursive: true });
}

// Expand ~ to the user's home directory
export function expandHomeDir(filePath: string): string {
  if (filePath.startsWith("~")) {
    return path.join(process.env.HOME || "", filePath.slice(1));
  }

  return filePath;
}

const PkgManagers = new Map<string, string>([
  ["bun", "bun i"],
  ["go", "go install"],
  ["npm", "npm i"],
  ["pnpm", "pnpm i"],
  ["yarn", "yarn install"],
]);

export function isRunnable(cmd: string): boolean {
  const shell = spawnSync(cmd);
  if (!shell) return false;
  return true;
}

export function isGoInstalled(): boolean {
  return isRunnable("go");
}

export function isTinyGoInstalled(): boolean {
  return isRunnable("tinygo");
}

export function isGitInstalled(): boolean {
  return isRunnable("git");
}

export function getGoVersion(): null | string {
  if (!isGoInstalled()) return null;
  const sh = execSync("go version").toString();
  return sh.split(" ")[2].slice(2);
}

export function getTinyGoVersion(): null | string {
  if (!isTinyGoInstalled()) return null;
  const sh = execSync("tinygo version").toString();
  return sh.split(" ")[2];
}

export async function cloneRepo(url: string, pth: string): Promise<boolean> {
  // https://github.com/hypermodeAI/tempalte-project/archive/refs/heads/main.zip instead
  const base_dir = path.dirname(pth);
  const temp_dir = path.join(base_dir, ".modus-temp");
  const folder_name = path.basename(pth);
  try {
    mkdirSync(base_dir, { recursive: true });
    if (existsSync(pth)) {
      const shell = execSync("git clone " + url + " .modus-temp", {
        stdio: "pipe",
        cwd: base_dir,
      });
      if (!shell) return false;
      cpSync(temp_dir, pth, { recursive: true, force: true });
      rmSync(temp_dir, { recursive: true });
      rmSync(path.join(pth, ".git"), { recursive: true });
      return true;
    }
    const shell = execSync("git clone " + url + " " + folder_name, {
      stdio: "pipe",
      cwd: base_dir,
    });
    if (!shell) return false;
    rmSync(path.join(pth, ".git"), { recursive: true });
    return true;
  } catch {
    return false;
  }
}

export function ask(question: string, rl: Interface, placeholder?: string): Promise<string> {
  return new Promise<string>((res, _) => {
    rl.question(question + (placeholder ? " " + placeholder + " " : ""), (answer) => {
      res(answer);
    });
  });
}

export function clearLine(): void {
  process.stdout.write(`\u001B[1A`);
  process.stdout.write("\u001B[2K");
  process.stdout.write("\u001B[0G");
}

export function getAvailablePackageManagers(): string[] {
  const pkgMgrs: string[] = [];
  if (isRunnable("npm")) pkgMgrs.push("NPM");
  if (isRunnable("yarn")) pkgMgrs.push("Yarn");
  if (isRunnable("pnpm")) pkgMgrs.push("PNPM");
  if (isRunnable("bun")) pkgMgrs.push("Bun");
  return pkgMgrs;
}

export function getLatestCLI(): string {
  // implement logic later
  return CLI_VERSION;
}
