# Copyright 2024 Hypermode, Inc.
# Licensed under the terms of the Apache License, Version 2.0
# See the LICENSE file that accompanied this code for further details.

# SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
# SPDX-License-Identifier: Apache-2.0

# Run with: iwr install.hypermode.com/modus.ps1 -UseBasicP | iex

# Config
$PACKAGE_ORG=""
$PACKAGE_NAME="test-modus-cli"
$INSTALL_DIR = "${env:MODUS_CLI}"
if (-not $INSTALL_DIR) {
    $INSTALL_DIR = "$HOME\.modus\cli"
}
$VERSION = "latest"

# Properties
$ARCH = $env:PROCESSOR_ARCHITECTURE
if ($ARCH -eq "AMD64") {
    $ARCH = "x64"
}
elseif ($ARCH -eq "ARM") {
    $ARCH = "arm"
}

# Functions
function Get-LatestVersion {
    $uri = if ($PACKAGE_ORG) {
        "https://registry.npmjs.org/$PACKAGE_ORG/$PACKAGE_NAME"
    } else {
        "https://registry.npmjs.org/$PACKAGE_NAME"
    }
    $response = Invoke-RestMethod -Uri $uri
    return $response."dist-tags".latest
}

function Install-Version {
    if ($VERSION -eq "latest") {
        $VERSION = Get-LatestVersion
    }

    Write-Output "${BOLD}${BLUE}Modus${RESET} Installer ${DIM}($VERSION)${RESET}`n"

    Install-Release

    echo "[4/5] Installed Modus CLI"
}

function Install-Release {
    Write-Host "[1/3] Fetching release from NPM"
    # Create a temporary file and get its directory
#     $tmpFile = New-TemporaryFile
# $tempFile = New-TemporaryFile
# $tempDir = [System.IO.Path]::ChangeExtension($tempFile.FullName, $null)
# New-Item -Path $tempDir -ItemType Directory -Force
# Remove-Item $tempFile.FullName
    $tmpDir = "$HOME\.modus-temp"
    $filename = "$PACKAGE_NAME-$VERSION.tgz"

    if (Test-Path $tmpDir) {
        Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    New-Item -Path $tmpDir -ItemType Directory -ErrorAction SilentlyContinue | Out-Null

    $downloadArchive = Download-ReleaseFromRepo $tmpDir
    Clear-Line
    Clear-Line

    Write-Host "[1/3] Fetched latest version ${DIM}$VERSION${RESET}"

    Write-Host "[2/3] Unpacking archive"
    
    tar -xf "$downloadArchive" -C "$tmpDir"
    Remove-Item -Path "$tmpDir/$filename" -Force -ErrorAction SilentlyContinue
    Clear-Line

    Write-Host "[2/3] Unpacked archive"
    if (Test-Path $INSTALL_DIR) {
        Remove-Item -Path $INSTALL_DIR -Recurse -Force -ErrorAction SilentlyContinue
    }

    New-Item -Path $INSTALL_DIR -ItemType Directory -ErrorAction SilentlyContinue | Out-Null
    Move-Item "$tmpDir/package/*" $INSTALL_DIR -ErrorAction SilentlyContinue
    Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue

    Write-Host "[3/3] Installed Modus CLI"

    Add-ToPath
    Restart-Shell
}

function Download-ReleaseFromRepo {
    param (
        [string]$tmpdir
    )
    $filename = "$PACKAGE_NAME-$VERSION.tgz"
    $downloadFile = Join-Path -Path $tmpdir -ChildPath $filename
    $archiveUrl = if ($PACKAGE_ORG) {
        "https://registry.npmjs.org/$PACKAGE_ORG/$PACKAGE_NAME/-/$filename"
    } else {
        "https://registry.npmjs.org/$PACKAGE_NAME/-/$filename"
    }
    cmd /c "curl --progress-bar --show-error --location --fail $archiveUrl --output $downloadFile"

    return $downloadFile
}

function Add-ToPath {
    $currentPath = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::User)

    if ($currentPath -notlike "*$INSTALL_DIR*") {
        $newPath = $currentPath + ";" + "$INSTALL_DIR\bin"
        [System.Environment]::SetEnvironmentVariable("Path", $newPath, [System.EnvironmentVariableTarget]::User)
        echo "[3/3] Added modus to PATH"
    } else {
        echo "[3/3] Modus already in PATH"
    }
}

function Restart-Shell {
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
}

function Clear-Line {
    Write-Host "${ESC}[F${ESC}[K" -NoNewLine
}

# ANSII codes
$ESC = [char]27
$BOLD = "${ESC}[1m"
$BLUE = "${ESC}[34;1m"
$DIM = "${ESC}[2m"
$RESET = "${ESC}[0m"

# This is the entry point

if ($env:MODUS_CLEAR_LINES) {
    for ($i = 1; $i -le $env:MODUS_CLEAR_LINES; $i++) {
        Clear-Line
    }
}

Install-Version
# Refresh Path
Restart-Shell
Write-Host "`nThe Modus CLI has been installed! " -NoNewLine
$(Write-Host ([System.char]::ConvertFromUtf32(127881)))
Write-Host "Run ${DIM}modus${RESET} to get started"
