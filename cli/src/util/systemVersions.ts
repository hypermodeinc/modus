import { execFile } from "./cp.js";
import { ModusHomeDir } from "../custom/globals.js";

export async function getGoVersion(): Promise<string | undefined> {
  try {
    const result = await execFile("go", ["version"], {
      shell: true,
      cwd: ModusHomeDir,
      env: process.env,
    });
    const parts = result.stdout.split(" ");
    const str = parts.length > 2 ? parts[2] : undefined;
    if (str?.startsWith("go")) {
      return str.slice(2);
    }
  } catch {}
}

export async function getTinyGoVersion(): Promise<string | undefined> {
  try {
    const result = await execFile("tinygo", ["version"], {
      shell: true,
      cwd: ModusHomeDir,
      env: process.env,
    });
    const parts = result.stdout.split(" ");
    return parts.length > 2 ? parts[2] : undefined;
  } catch {}
}

export async function getNPMVersion(): Promise<string | undefined> {
  try {
    const result = await execFile("npm", ["--version"], { shell: true });
    const parts = result.stdout.split(" ");
    return parts.length > 0 ? parts[0].trim() : undefined;
  } catch {}
}
