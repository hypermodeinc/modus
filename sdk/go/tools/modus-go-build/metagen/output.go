/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package metagen

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/metadata"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/utils"
)

func LogToConsole(meta *metadata.Metadata) {

	// FORCE_COLOR is set by Modus CLI
	forceColor := os.Getenv("FORCE_COLOR")
	if forceColor != "" && forceColor != "0" {
		color.NoColor = false
	}
	w := color.Output

	writeHeader(w, "Metadata:")
	writeTable(w, [][]string{
		{"Plugin Name", meta.Plugin},
		{"Go Module", meta.Module},
		{"Modus SDK", meta.SDK},
		{"Build ID", meta.BuildId},
		{"Build Timestamp", meta.BuildTime},
		{"Git Repository", meta.GitRepo},
		{"Git Commit", meta.GitCommit},
	})
	fmt.Fprintln(w)

	if len(meta.FnExports) > 0 {
		writeHeader(w, "Functions:")
		for _, k := range meta.FnExports.SortedKeys() {
			fn := meta.FnExports[k]
			writeItem(w, fn.String(meta))
		}
		fmt.Fprintln(w)
	}

	types := make([]string, 0, len(meta.Types))
	for _, k := range meta.Types.SortedKeys(meta.Module) {
		t := meta.Types[k]
		if len(t.Fields) > 0 && strings.HasPrefix(k, meta.Module) {
			types = append(types, k)
		}
	}

	if len(types) > 0 {
		writeHeader(w, "Custom Types:")
		for _, t := range types {
			writeItem(w, meta.Types[t].String(meta))
		}
		fmt.Fprintln(w)
	}

	if utils.IsDebugModeEnabled() {
		writeHeader(w, "Metadata JSON:")
		metaJson, _ := utils.JsonSerialize(meta, true)
		fmt.Fprintln(w, string(metaJson))
		fmt.Fprintln(w)
	}

}

func writeHeader(w io.Writer, text string) {
	color.Set(color.FgBlue, color.Bold)
	fmt.Fprintln(w, text)
	color.Unset()
}

func writeItem(w io.Writer, text string) {
	color.Set(color.FgCyan)
	fmt.Fprint(w, "  "+text)
	color.Unset()
	fmt.Fprintln(w)
}

func writeTable(w io.Writer, rows [][]string) {
	pad := make([]int, len(rows))
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > pad[i] {
				pad[i] = len(cell)
			}
		}
	}

	for _, row := range rows {
		if len(row) != 2 || len(row[0]) == 0 || len(row[1]) == 0 {
			continue
		}

		padding := strings.Repeat(" ", pad[0]-len(row[0]))

		fmt.Fprint(w, "  ")
		color.Set(color.FgCyan)
		fmt.Fprintf(w, "%s:%s ", row[0], padding)
		color.Set(color.FgBlue)
		fmt.Fprint(w, row[1])
		color.Unset()
		fmt.Fprintln(w)
	}
}
