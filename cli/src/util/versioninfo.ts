/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import semver from "semver";
import path from "node:path";
import * as fs from "./fs.js";
import * as globals from "../custom/globals.js";
import { getGitHubApiHeaders } from "./index.js";

export function getSdkPath(sdk: globals.SDK, version: string): string {
  return path.join(globals.ModusHomeDir, "sdk", sdk.toLowerCase(), version);
}

export function getRuntimePath(version: string): string {
  return path.join(globals.ModusHomeDir, "runtime", version);
}

export function isPrerelease(version: string): boolean {
  if (version.startsWith("v")) {
    version = version.slice(1);
  }
  return !!semver.prerelease(version);
}

export async function fetchFromModusLatestNoPrerelease(): Promise<any> {
  const response = await fetch(`https://releases.hypermode.com/modus-latest.json`, {});
  if (!response.ok) {
    throw new Error(`Error fetching latest SDK version: ${response.statusText}`);
  }

  return await response.json();
}

export async function fetchFromModusAllNoPrerelease(): Promise<any> {
  const response = await fetch(`https://releases.hypermode.com/modus-all.json`, {});
  if (!response.ok) {
    throw new Error(`Error fetching all SDK versions: ${response.statusText}`);
  }

  return await response.json();
}

export async function fetchItemVersionsFromModusAllNoPrerelease(item: string): Promise<string[]> {
  const data = await fetchFromModusAllNoPrerelease();

  let versions: string[];
  switch (item) {
    case "sdk/assemblyscript/":
      versions = data["sdk/assemblyscript"];
      break;
    case "sdk/go/":
      versions = data["sdk/go"];
      break;
    case "runtime/":
      versions = data["runtime"];
      break;
    case "cli/":
      versions = data["cli"];
      break;
    default:
      throw new Error("Not a valid item in releases");
  }

  return versions.map((version) => `v${version}`);
}

export async function getLatestSdkVersionNoPrerelease(sdk: globals.SDK): Promise<string | undefined> {
  const data = await fetchFromModusLatestNoPrerelease();
  if (sdk === globals.SDK.AssemblyScript) {
    return data["sdk/assemblyscript"];
  } else if (sdk === globals.SDK.Go) {
    return "v" + data["sdk/go"];
  }
}

export async function getLatestRuntimeVersionNoPrerelease(): Promise<string | undefined> {
  const data = await fetchFromModusLatestNoPrerelease();
  return "v" + data["runtime"];
}

export async function getLatestCliVersionNoPrerelease(): Promise<string | undefined> {
  const data = await fetchFromModusLatestNoPrerelease();
  return "v" + data["cli"];
}

export async function getLatestSdkVersion(sdk: globals.SDK, includePrerelease: boolean): Promise<string | undefined> {
  if (!includePrerelease) {
    return await getLatestSdkVersionNoPrerelease(sdk);
  }
  return await getLatestVersion(globals.GitHubOwner, globals.GitHubRepo, globals.GetSdkTagPrefix(sdk), includePrerelease);
}

export async function getLatestRuntimeVersion(includePrerelease: boolean): Promise<string | undefined> {
  if (!includePrerelease) {
    return await getLatestRuntimeVersionNoPrerelease();
  }
  return await getLatestVersion(globals.GitHubOwner, globals.GitHubRepo, globals.GitHubRuntimeTagPrefix, includePrerelease);
}

async function getLatestVersion(owner: string, repo: string, prefix: string, includePrerelease: boolean): Promise<string | undefined> {
  try {
    let tag = await findLatestReleaseTag(owner, repo, prefix, includePrerelease);
    if (!tag && !includePrerelease) {
      // If no stable release was found, look for a prerelease
      tag = await findLatestReleaseTag(owner, repo, prefix, true);
    }
    if (tag) {
      return tag.slice(prefix.length);
    }
  } catch (e) {
    console.error(e);
  }
}

export async function getAllSdkVersions(sdk: globals.SDK, includePrerelease: boolean): Promise<string[]> {
  return await getAllVersions(globals.GitHubOwner, globals.GitHubRepo, globals.GetSdkTagPrefix(sdk), includePrerelease);
}

export async function getAllRuntimeVersions(includePrerelease: boolean): Promise<string[]> {
  return await getAllVersions(globals.GitHubOwner, globals.GitHubRepo, globals.GitHubRuntimeTagPrefix, includePrerelease);
}

async function getAllVersions(owner: string, repo: string, prefix: string, includePrerelease: boolean): Promise<string[]> {
  try {
    if (!includePrerelease) {
      return await fetchItemVersionsFromModusAllNoPrerelease(prefix);
    }
    let tags = await getAllReleaseTags(owner, repo, prefix, includePrerelease);
    if (tags.length === 0 && !includePrerelease) {
      // If no stable release was found, look for prereleases
      tags = await getAllReleaseTags(owner, repo, prefix, true);
    }
    let versions = tags.map((tag) => {
      let version = tag.slice(prefix.length);
      if (version.startsWith("v")) {
        version = version.slice(1);
      }
      return version;
    });
    versions = semver.rsort(versions);
    return versions.map((v) => "v" + v);
  } catch (e) {
    console.error(e);
    return [];
  }
}

export async function sdkReleaseExists(sdk: globals.SDK, version: string): Promise<boolean> {
  const prefix = globals.GetSdkTagPrefix(sdk);
  return releaseExists(globals.GitHubOwner, globals.GitHubRepo, `${prefix}${version}`);
}

export async function runtimeReleaseExists(version: string): Promise<boolean> {
  return releaseExists(globals.GitHubOwner, globals.GitHubRepo, `${globals.GitHubRuntimeTagPrefix}${version}`);
}

export async function sdkVersionIsInstalled(sdk: globals.SDK, version: string): Promise<boolean> {
  return await fs.exists(getSdkPath(sdk, version));
}

export async function runtimeVersionIsInstalled(version: string): Promise<boolean> {
  return await fs.exists(getRuntimePath(version));
}

export async function getLatestInstalledSdkVersion(sdk: globals.SDK, includePrerelease: boolean): Promise<string | undefined> {
  const dir = path.join(globals.ModusHomeDir, "sdk", sdk.toLowerCase());
  const versions = await getInstalledVersions(dir, includePrerelease);
  return versions.length > 0 ? versions[0] : undefined;
}

export async function getLatestInstalledRuntimeVersion(includePrerelease: boolean): Promise<string | undefined> {
  const dir = path.join(globals.ModusHomeDir, "runtime");
  const versions = await getInstalledVersions(dir, includePrerelease);
  return versions.length > 0 ? versions[0] : undefined;
}

export async function getInstalledSdkVersions(sdk: globals.SDK): Promise<string[]> {
  const dir = path.join(globals.ModusHomeDir, "sdk", sdk.toLowerCase());
  return await getInstalledVersions(dir, true);
}

export async function getInstalledRuntimeVersions(): Promise<string[]> {
  const dir = path.join(globals.ModusHomeDir, "runtime");
  return await getInstalledVersions(dir, true);
}

async function getInstalledVersions(dir: string, includePrerelease: boolean): Promise<string[]> {
  if (await fs.exists(dir)) {
    const entries = await fs.readdir(dir, { withFileTypes: true });
    let versions = entries.filter((e) => e.isDirectory() && e.name.startsWith("v")).map((e) => e.name.slice(1));

    // If at least one stable release is found, only return stable releases.
    // Otherwise, allow prereleases, regardless of the includePrerelease flag.
    if (!includePrerelease) {
      const stableVersions = versions.filter((v) => !semver.prerelease(v));
      if (stableVersions.length > 0) {
        versions = stableVersions;
      }
    }

    return semver.rsort(versions).map((v) => "v" + v);
  }
  return [];
}

export async function findCompatibleInstalledRuntimeVersion(sdk: globals.SDK, sdkVersion: string, includePrerelease: boolean): Promise<string | undefined> {
  const infoFile = path.join(getSdkPath(sdk, sdkVersion), "sdk.json");
  if (!(await fs.exists(infoFile))) {
    throw new Error(`SDK info file not found: ${infoFile}`);
  }

  const info = JSON.parse(await fs.readFile(infoFile, "utf8"));
  let constraint: string = info.runtime.version;
  if (constraint.startsWith("v")) {
    constraint = constraint.slice(1);
  }

  if (semver.valid(constraint)) {
    // exact version specified - check if it is installed
    const version = "v" + constraint;
    if (await runtimeVersionIsInstalled(version)) {
      return version;
    }
  } else if (semver.validRange(constraint)) {
    // range specified - find the latest installed version that satisfies it
    const versions = await getInstalledRuntimeVersions();
    let compatibleVersions = versions.filter((v) => semver.satisfies(v.slice(1), constraint, { includePrerelease }));
    if (compatibleVersions.length == 0 && !includePrerelease) {
      // If no stable release was found, look for a prerelease
      compatibleVersions = versions.filter((v) => semver.satisfies(v.slice(1), constraint, { includePrerelease: true }));
    }
    if (compatibleVersions.length > 0) {
      return compatibleVersions[0];
    }
  } else {
    throw new Error(`Invalid runtime version: ${constraint}`);
  }
}

export async function findLatestCompatibleRuntimeVersion(sdk: globals.SDK, sdkVersion: string, includePrerelease: boolean): Promise<string | undefined> {
  const infoFile = path.join(getSdkPath(sdk, sdkVersion), "sdk.json");
  if (!(await fs.exists(infoFile))) {
    throw new Error(`SDK info file not found: ${infoFile}`);
  }

  const info = JSON.parse(await fs.readFile(infoFile, "utf8"));
  let constraint = info.runtime.version;
  if (constraint.startsWith("v")) {
    constraint = constraint.slice(1);
  }

  if (semver.valid(constraint) !== null) {
    // exact version specified - check if it exists
    const version = "v" + constraint;
    if (await runtimeReleaseExists(version)) {
      return version;
    }
  } else if (semver.validRange(constraint, { includePrerelease }) !== null) {
    // range specified - find the latest released version that satisfies it
    const versions = await getAllRuntimeVersions(includePrerelease);
    let compatibleVersions = versions.filter((v) => semver.satisfies(v.slice(1), constraint, { includePrerelease }));
    if (compatibleVersions.length == 0 && !includePrerelease) {
      // If no stable release was found, look for a prerelease
      compatibleVersions = versions.filter((v) => semver.satisfies(v.slice(1), constraint, { includePrerelease: true }));
    }
    if (compatibleVersions.length > 0) {
      return compatibleVersions[0];
    }
  }
}

const headers = getGitHubApiHeaders();

async function releaseExists(owner: string, repo: string, tag: string): Promise<boolean> {
  const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/releases/tags/${encodeURIComponent(tag)}`, { headers });
  return response.ok;
}

async function findLatestReleaseTag(owner: string, repo: string, prefix: string, includePrerelease: boolean): Promise<string | undefined> {
  let page = 1;
  while (true) {
    const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/releases?page=${page}`, { headers });

    if (!response.ok) {
      throw new Error(`Error fetching releases: ${response.statusText}`);
    }

    const releases = await response.json();
    if (releases.length === 0) {
      return;
    }

    for (const release of releases) {
      if (!includePrerelease && release.prerelease) {
        continue;
      }

      if (prefix && !release.tag_name.startsWith(prefix)) {
        continue;
      }

      return release.tag_name;
    }

    page++;
  }
}

async function getAllReleaseTags(owner: string, repo: string, prefix: string, includePrerelease: boolean): Promise<string[]> {
  const results: string[] = [];

  let page = 1;
  while (true) {
    const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/releases?per_page=100&page=${page}`, { headers });

    if (!response.ok) {
      throw new Error(`Error fetching releases: ${response.statusText}`);
    }

    const releases = await response.json();
    if (releases.length === 0) {
      return results;
    }

    for (const release of releases) {
      if (!includePrerelease && release.prerelease) {
        continue;
      }

      if (prefix && !release.tag_name.startsWith(prefix)) {
        continue;
      }

      results.push(release.tag_name);
    }

    page++;
  }
}
