# Copyright 2024 Hypermode, Inc.
# Licensed under the terms of the Apache License, Version 2.0
# See the LICENSE file that accompanied this code for further details.
#
# SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
# SPDX-License-Identifier: Apache-2.0

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

    Write-Output "${BOLD}${BLUE}Modus${RESET} Installer ($VERSION)`n"

    Install-Release
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

    Write-Host "[2/3] Unpacking archive"
    
    tar -xf "$downloadArchive" -C "$tmpDir"
    Remove-Item -PAth "$tmpDir/modus-$VERSION-win32-$ARCH.zip" -Force

    if (Test-Path $INSTALL_DIR) {
        Remove-Item -Path $INSTALL_DIR -Recurse -Force
    }

    New-Item -Path $INSTALL_DIR -ItemType Directory | Out-Null
    Move-Item "$tmpDir/modus/*" $INSTALL_DIR
    Remove-Item -Path $tmpDir -Recurse -Force

    Write-Host "[3/3] Installed Modus CLI"

    # Add it to path
    # Uh, finalize, clean up, restart terminal
}

function Download-ReleaseFromRepo {
    param (
        [string]$tmpdir
    )
    $filename = "modus-$VERSION-win32-$ARCH.zip"
    $downloadFile = Join-Path -Path $tmpdir -ChildPath $filename
    $archiveUrl = "https://github.com/$GIT_REPO/releases/download/$VERSION/$filename"

    cmd /c "curl --progress-bar --show-error --location --fail $archiveUrl --output $downloadFile"

    return $downloadFile
}


$ESC = [char]27

$BOLD = "${ESC}[1m"
$BLUE = "${ESC}[34;1m"
$DIM = "${ESC}[2m"
$RESET = "${ESC}[0m"
# This is the entry point
Install-Version
# Update-Profile $INSTALL_DIR
# Write-Host "The Modus CLI has been installed! 🎉"
# Write-Host "Run 'modus' to get started."
