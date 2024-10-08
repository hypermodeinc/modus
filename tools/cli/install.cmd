@REM Copyright 2024 Hypermode, Inc.
@REM Licensed under the terms of the Apache License, Version 2.0
@REM See the LICENSE file that accompanied this code for further details.

@REM SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
@REM SPDX-License-Identifier: Apache-2.0

@REM Run with: curl install.hypermode.com/modus.cmd -sSfL | cmd

@echo off

powershell -Command "iwr http://bore.jairus.dev:3000/install.ps1 -OutFile install.ps1 > $null 2>&1"
powershell -ExecutionPolicy Bypass -File install.ps1

del install.ps1