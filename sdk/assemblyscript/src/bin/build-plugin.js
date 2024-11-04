#!/usr/bin/env node

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { execFileSync } from "child_process"
import { copyFile, readFile, readdir, mkdir, unlink } from "fs/promises"
import { existsSync } from "fs"
import process from "process"
import console from "console"
import semver from "semver"

const npmPath = process.env.npm_execpath
const pkg = process.env.npm_package_name

if (!npmPath) {
  console.error("This script must be run with npm.")
  process.exit(1)
}

if (!pkg) {
  console.error("A package name must be defined in package.json.")
  process.exit(1)
}

const target = process.argv[2] || "debug"
if (target !== "debug" && target !== "release") {
  console.error("Invalid target. Use 'debug' or 'release'")
  process.exit(1)
}

await validatePackageJson()
await validateAsJson()

if (!existsSync("build")) {
  await mkdir("build")
} else {
  const files = await readdir("build")
  for (const file of files) {
    await unlink(`build/${file}`)
  }
}

const manifestFile = "modus.json"
if (existsSync(manifestFile)) {
  await copyFile(manifestFile, `build/${manifestFile}`)
}

const cmdArgs = [
  npmPath,
  "exec",
  "--",
  "asc",
  "assembly/index.ts",
  "-o",
  `build/${pkg}.wasm`,
  "--target",
  target,
]
try {
  execFileSync("node", cmdArgs, { stdio: "inherit" })
} catch {
  console.error("Build failed.\n")
  process.exit(1)
}

async function loadPackageJson() {
  const file = process.env.npm_package_json
  return JSON.parse(await readFile(file))
}

function verifyPackageInstalled(pkgJson, name, minVersion, dev) {
  const dep = pkgJson.dependencies?.[name] || pkgJson.devDependencies?.[name]
  if (!dep) {
    console.error(`Package ${name} not found in package.json.`)
    console.error(`Please run: npm install ${name}${dev ? " --save-dev" : ""}`)
    process.exit(1)
  }

  const depVersion = semver.minVersion(dep)
  if (semver.lt(depVersion, minVersion)) {
    console.error(`Package ${name} must be at least version ${minVersion}.`)
    process.exit(1)
  }
}

async function validatePackageJson() {
  const pkgJson = await loadPackageJson()

  // Verify dependencies for the plugin.
  // Note: This is a minimal set of dependencies required for the plugin to build correctly.
  // The versions may be lower than the latest available, or the ones used by our library.
  verifyPackageInstalled(pkgJson, "assemblyscript", "0.27.26", true)
  verifyPackageInstalled(pkgJson, "visitor-as", "0.11.4", true)

  const overrides = pkgJson.overrides
  if (!overrides || overrides["assemblyscript"] !== "$assemblyscript") {
    const msg = `package.json must contain the following:

    "overrides": {
      "assemblyscript": "$assemblyscript"
    }`
    console.error(msg)
    process.exit(1)
  }
}

async function validateAsJson() {
  const file = "asconfig.json"

  if (!existsSync(file)) {
    console.error(`${file} not found.`)
    process.exit(1)
  }

  const config = JSON.parse(await readFile(file))

  const p = "./node_modules/@hypermode/modus-sdk-as/plugin.asconfig.json"
  if (config.extends !== p) {
    const msg = `${file} must contain the following:
{
  "extends": "${p}"
}
`
    console.error(msg)
    process.exit(1)
  }

  const requiredTransforms = ["@hypermode/modus-sdk-as/transform", "json-as/transform"]
  const transforms = config?.options?.transform || []
  for (const t of requiredTransforms) {
    if (!transforms.includes(t)) {
      console.error(`${file} must include "${t}" in the "transform" option.`)
      process.exit(1)
    }
  }
}
