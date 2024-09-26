import { execSync, spawnSync } from "node:child_process";
import { existsSync, mkdirSync, rmdirSync } from "node:fs";
import path from "node:path";
import { createInterface, Interface } from "node:readline";

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
    ["yarn", "yarn install"]
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
    return isRunnable("git")
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

export function cloneRepo(url: string, dstFolder?: string): boolean {
    const shell = execSync("git clone " + url + " " + dstFolder, { stdio: "pipe" });
    if (!shell) return false;
    return true;
}

export function installDeps(manager: string, name: string): boolean {
    // this is essential to sanitize user input
    const cmd = PkgManagers.get(manager.toLowerCase());
    if (!cmd) return false;
    const shell = execSync(cmd, { cwd: path.join(process.cwd(), name), stdio: "ignore" });

    return true;
}

export function ask(question: string, face: Interface): Promise<string> {
    return new Promise<string>((res, _) => {
        face.question(question, (answer) => {
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
    return "0.0.1";
}