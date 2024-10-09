# Copyright 2024 Hypermode, Inc.
# Licensed under the terms of the Apache License, Version 2.0
# See the LICENSE file that accompanied this code for further details.

# SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
# SPDX-License-Identifier: Apache-2.0

# Run with: iwr install.hypermode.com/modus.ps1 -UseBasicP | iex

# Config
$GIT_REPO = "JairusSW/modus-cli"
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
function Get-LatestRelease {
    $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$GIT_REPO/releases/latest"
    return $response.tag_name
}

function Install-Version {
    if ($VERSION -eq "latest") {
        $VERSION = Get-LatestRelease
    }

    Write-Output "${BOLD}${BLUE}Modus${RESET} Installer ${DIM}($VERSION)${RESET}`n"

    Install-Release

    echo "[4/5] Installed Modus CLI"
}
function Install-Release {
    Write-Host "[1/3] Fetching archive for Windows $ARCH"
    # Create a temporary file and get its directory
#     $tmpFile = New-TemporaryFile
# $tempFile = New-TemporaryFile
# $tempDir = [System.IO.Path]::ChangeExtension($tempFile.FullName, $null)
# New-Item -Path $tempDir -ItemType Directory -Force
# Remove-Item $tempFile.FullName
    $tmpDir = "$HOME\.modus-temp"

    if (Test-Path $tmpDir) {
        Remove-Item -Path $tmpDir -Recurse -Force
    }
    New-Item -Path $tmpDir -ItemType Directory | Out-Null

    $downloadArchive = Download-ReleaseFromRepo $tmpDir
    Clear-Line
    Clear-Line

    Write-Host "[1/3] Fetched archive for Windows $ARCH"

    Write-Host "[2/3] Unpacking archive"
    
    tar -xf "$downloadArchive" -C "$tmpDir"
    Remove-Item -Path "$tmpDir/modus-$VERSION-win32-$ARCH.zip" -Force
    Clear-Line

    Write-Host "[2/3] Unpacked archive"
    if (Test-Path $INSTALL_DIR) {
        Remove-Item -Path $INSTALL_DIR -Recurse -Force
    }

    New-Item -Path $INSTALL_DIR -ItemType Directory | Out-Null
    Move-Item "$tmpDir/*" $INSTALL_DIR
    Remove-Item -Path $tmpDir -Recurse -Force

    Write-Host "[3/3] Installed Modus CLI"

    Add-ToPath
    Restart-Shell
}

function Download-ReleaseFromRepo {
    param (
        [string]$tmpdir
    )
    $filename = "modus-cli-$VERSION.zip"
    $downloadFile = Join-Path -Path $tmpdir -ChildPath $filename
    $archiveUrl = "https://github.com/$GIT_REPO/releases/download/$VERSION/$filename"

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

function Clear-Banner {
    # Define the path to the install.cmd file
    $installCmdPath = "./install.cmd"
    
    # Check if the install.cmd file exists
    if (Test-Path $installCmdPath) {
        # If it exists, clear the line 3 times
        for ($i = 0; $i -lt 7; $i++) {
            Clear-Line
        }
    } else {
        Write-Host "install.cmd does not exist."
    }
}

# ANSII codes
$ESC = [char]27
$BOLD = "${ESC}[1m"
$BLUE = "${ESC}[34;1m"
$DIM = "${ESC}[2m"
$RESET = "${ESC}[0m"

# This is the entry point
Clear-Banner
Install-Version
Write-Host "`nThe Modus CLI has been installed! " -NoNewLine
$(Write-Host ([System.char]::ConvertFromUtf32(127881)))
Write-Host "Run ${DIM}modus${RESET} to get started"
