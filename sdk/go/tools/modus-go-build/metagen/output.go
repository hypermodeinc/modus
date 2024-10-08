/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package metagen

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/metadata"
	"github.com/hypermodeinc/modus/sdk/go/tools/modus-go-build/utils"
	"github.com/mattn/go-isatty"
)

type outputMode int

const (
	outputText     outputMode = 1
	outputMarkdown outputMode = 2
)

//go:embed logo.txt
var logo string

func WriteLogo() {
	w := os.Stdout
	mode := getOutputMode(w)

	fmt.Fprintln(w)
	writeAsciiLogo(w, mode)
}

func LogToConsole(meta *metadata.Metadata) {
	w := os.Stdout
	mode := getOutputMode(w)

	writeHeader(w, mode, "Metadata:")
	writeTable(w, mode, [][]string{
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
		writeHeader(w, mode, "Functions:")
		for _, k := range meta.FnExports.SortedKeys() {
			fn := meta.FnExports[k]
			writeItem(w, mode, fn.String(meta))
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
		writeHeader(w, mode, "Custom Types:")
		for _, t := range types {
			writeItem(w, mode, meta.Types[t].String(meta))
		}
		fmt.Fprintln(w)
	}

	if utils.IsDebugModeEnabled() {
		writeHeader(w, mode, "Metadata JSON:")
		metaJson, _ := utils.JsonSerialize(meta, true)
		fmt.Fprintln(w, string(metaJson))
		fmt.Fprintln(w)
	}

}

func getOutputMode(f *os.File) outputMode {
	color.NoColor = !isatty.IsTerminal(f.Fd())
	color.Output = f
	return outputText
}

func writeAsciiLogo(w io.Writer, mode outputMode) {
	switch mode {
	case outputText:
		color.Set(color.FgBlue)
		fmt.Fprintln(w, logo)
		color.Unset()
	case outputMarkdown:
		fmt.Fprintln(w, "```")
		fmt.Fprintln(w, logo)
		fmt.Fprintln(w, "```")
	}
}

func writeHeader(w io.Writer, mode outputMode, text string) {
	switch mode {
	case outputText:
		color.Set(color.FgBlue, color.Bold)
		fmt.Fprintln(w, text)
		color.Unset()
	case outputMarkdown:
		fmt.Fprintf(w, "### %s\n", text)
	}
}

func writeItem(w io.Writer, mode outputMode, text string) {
	switch mode {
	case outputText:
		color.Set(color.FgCyan)
		fmt.Fprintf(w, "  %s\n", text)
		color.Unset()
	case outputMarkdown:
		fmt.Fprintf(w, "  - %s\n", text)
	}
}

func writeTable(w io.Writer, mode outputMode, rows [][]string) {
	pad := make([]int, len(rows))
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > pad[i] {
				pad[i] = len(cell)
			}
		}
	}

	if mode == outputMarkdown {
		fmt.Fprintf(w, "| %s | %s |\n", strings.Repeat(" ", pad[0]), strings.Repeat(" ", pad[1]))
		fmt.Fprintf(w, "| %s | %s |\n", strings.Repeat("-", pad[0]), strings.Repeat("-", pad[1]))
	}

	for _, row := range rows {
		if len(row) != 2 || len(row[0]) == 0 || len(row[1]) == 0 {
			continue
		}

		pad0 := strings.Repeat(" ", pad[0]-len(row[0]))
		pad1 := strings.Repeat(" ", pad[1]-len(row[1]))

		switch mode {
		case outputText:
			fmt.Fprint(w, "  ")
			color.Set(color.FgCyan)
			fmt.Fprintf(w, "%s:%s ", row[0], pad0)
			color.Set(color.FgBlue)
			fmt.Fprint(w, row[1])
			color.Unset()
			fmt.Fprint(w, "\n")

		case outputMarkdown:
			fmt.Fprintf(w, "| %s%s | %s%s |\n", row[0], pad0, row[1], pad1)
		}
	}
}
