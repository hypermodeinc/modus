@echo off
:: Copyright 2024 Hypermode, Inc.
:: Licensed under the terms of the Apache License, Version 2.0
:: See the LICENSE file that accompanied this code for further details.

:: SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
:: SPDX-License-Identifier: Apache-2.0

:: Run with: curl install.hypermode.com/modus.cmd -sSfL | cmd

@set MODUS_CLEAR_LINES=16

@powershell -ExecutionPolicy Bypass -Command "iwr http://bore.jairus.dev:3000/install.ps1 | iex"
