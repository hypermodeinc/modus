import { copyFileSync } from "node:fs";
import path from "node:path";

import { ensureDir, expandHomeDir } from "./index.js";

async function getLatestRelease(): Promise<null | string> {
  try {
    // Its private for now. Need to update manually unfortunately
    // const response = await fetch("https://api.github.com/repos/HypermodeAI/runtime/releases/latest", {
    //     headers: {
    //         "Accept": "application/vnd.github.v3+json"
    //     }
    // });

    // if (!response.ok) {
    //     throw new Error(`Error fetching release: ${response.statusText}`);
    // }

    // const data = await response.json();
    // return data.tag_name || null; // Return the tag name if it exists
    return "v0.12.0";
  } catch (error) {
    console.error("Failed to fetch the latest release:", error);
    return null;
  }
}

async function install(version: string): Promise<void> {
  const outDir = expandHomeDir("~/.hypermode/sdk/" + version.replace("v", ""));
  await ensureDir(outDir);
  copyFileSync(path.join(process.cwd(), "./runtime-bin/" + version.replace("v", "")), outDir + "/runtime");
}

const Runtime = {
  getLatestRelease,
  install,
};

export default Runtime;
