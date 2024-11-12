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
import * as http from "./http.js";
import * as fs from "./fs.js";
import * as globals from "../custom/globals.js";

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

const versionData: {
  latest?: { [key: string]: string };
  preview?: { [key: string]: string };
  all?: { [key: string]: string[] };
  allPreview?: { [key: string]: string[] };
} = {};

export async function fetchModusLatest() {
  if (versionData.latest) return versionData.latest;

  const response = await http.get(`https://releases.hypermode.com/modus-latest.json`);
  if (!response.ok) {
    throw new Error(`Error fetching latest SDK versions: ${response.statusText}`);
  }

  const data = (await response.json()) as { [key: string]: string };
  if (!data) {
    throw new Error(`Error fetching latest SDK versions: response was empty`);
  }

  versionData.latest = data;
  return data;
}

export async function fetchModusPreview() {
  if (versionData.preview) return versionData.preview;

  const response = await http.get(`https://releases.hypermode.com/modus-preview.json`);
  if (!response.ok) {
    throw new Error(`Error fetching latest preview SDK versions: ${response.statusText}`);
  }

  const data = (await response.json()) as { [key: string]: string };
  if (!data) {
    throw new Error(`Error fetching latest preview SDK versions: response was empty`);
  }

  versionData.preview = data;
  return data;
}

export async function fetchModusAll() {
  if (versionData.all) return versionData.all;

  const response = await http.get(`https://releases.hypermode.com/modus-all.json`);
  if (!response.ok) {
    throw new Error(`Error fetching all SDK versions: ${response.statusText}`);
  }

  const data = (await response.json()) as { [key: string]: string[] };
  if (!data) {
    throw new Error(`Error fetching all SDK versions: response was empty`);
  }

  versionData.all = data;
  return data;
}

export async function fetchModusPreviewAll() {
  if (versionData.allPreview) return versionData.allPreview;

  const response = await http.get(`https://releases.hypermode.com/modus-preview-all.json`);
  if (!response.ok) {
    throw new Error(`Error fetching all preview SDK versions: ${response.statusText}`);
  }

  const data = (await response.json()) as { [key: string]: string[] };
  if (!data) {
    throw new Error(`Error fetching all preview SDK versions: response was empty`);
  }

  versionData.allPreview = data;
  return data;
}

export async function fetchItemVersionsFromModusAll(item: string): Promise<string[]> {
  const data = await fetchModusAll();

  if (item.endsWith("/")) {
    item = item.slice(0, -1);
  }

  const versions = data[item];
  if (!versions) {
    throw new Error("Not a valid item in releases");
  }

  return versions;
}

export async function fetchItemVersionsFromModusPreviewAll(item: string): Promise<string[]> {
  const data = await fetchModusPreviewAll();

  if (item.endsWith("/")) {
    item = item.slice(0, -1);
  }

  const versions = data[item];
  if (!versions) {
    throw new Error("Not a valid item in releases");
  }

  return versions;
}

export async function getLatestSdkVersion(sdk: globals.SDK, includePrerelease: boolean): Promise<string | undefined> {
  const data = includePrerelease ? await fetchModusPreview() : await fetchModusLatest();
  const version = data["sdk/" + sdk.toLowerCase()];
  return version ? version : undefined;
}

export async function getLatestRuntimeVersion(includePrerelease: boolean): Promise<string | undefined> {
  const data = includePrerelease ? await fetchModusPreview() : await fetchModusLatest();
  const version = data["runtime"];
  return version ? version : undefined;
}

export async function getLatestCliVersion(includePrerelease: boolean): Promise<string | undefined> {
  const data = includePrerelease ? await fetchModusPreview() : await fetchModusLatest();
  const version = data["cli"];
  return version ? version : undefined;
}

export async function getAllSdkVersions(sdk: globals.SDK, includePrerelease: boolean): Promise<string[]> {
  return await getAllVersions(globals.GetSdkTagPrefix(sdk), includePrerelease);
}

export async function getAllRuntimeVersions(includePrerelease: boolean): Promise<string[]> {
  return await getAllVersions(globals.GitHubRuntimeTagPrefix, includePrerelease);
}

async function getAllVersions(prefix: string, includePrerelease: boolean): Promise<string[]> {
  try {
    return includePrerelease ? await fetchItemVersionsFromModusPreviewAll(prefix) : await fetchItemVersionsFromModusAll(prefix);
  } catch (e) {
    console.error(e);
    return [];
  }
}

export async function sdkReleaseExists(sdk: globals.SDK, version: string): Promise<boolean> {
  const versions = await getAllSdkVersions(sdk, true);
  return versions.includes(version);
}

export async function runtimeReleaseExists(version: string): Promise<boolean> {
  const versions = await getAllRuntimeVersions(true);
  return versions.includes(version);
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
