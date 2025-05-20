/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/rs/zerolog"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func useEmbeddedPostgres() bool {
	s := os.Getenv("MODUS_DB_USE_EMBEDDED_POSTGRES")
	if s != "" {
		if value, err := strconv.ParseBool(s); err == nil {
			return value
		}
	}
	return runtime.GOOS == "windows"
}

var _embeddedPostgresDB *embeddedpostgres.EmbeddedPostgres

func getEmbeddedPostgresDataDir(ctx context.Context) string {
	var dataDir string
	appPath := app.Config().AppPath()
	if filepath.Base(appPath) == "build" {
		// this keeps the data directory outside of the build directory
		dataDir = filepath.Join(appPath, "..", ".postgres")
		addToGitIgnore(ctx, filepath.Dir(appPath), ".postgres/")
	} else {
		dataDir = filepath.Join(appPath, ".postgres")
	}
	return dataDir
}

func prepareEmbeddedPostgresDB(ctx context.Context) error {

	dataDir := getEmbeddedPostgresDataDir(ctx)

	if err := shutdownPreviousEmbeddedPostgresDB(ctx, dataDir); err != nil {
		return fmt.Errorf("error shutting down previous embedded postgres instance: %w", err)
	}

	port, err := findAvailablePort(5432, 5499)
	if err != nil {
		return err
	}

	logger.Info(ctx).Msg("Preparing embedded PostgreSQL database.  The db instance will log its output next:")

	// Note: releases come from here:
	// https://github.com/zonkyio/embedded-postgres-binaries/releases

	dbname := "modusdb"
	cfg := embeddedpostgres.DefaultConfig().
		Port(port).
		Database(dbname).
		DataPath(dataDir).
		Version("17.4.0").
		OwnProcessGroup(true).
		BinariesPath(getEmbeddedPostgresBinariesPath())

	db := embeddedpostgres.NewDatabase(cfg)
	if err := db.Start(); err != nil {
		return fmt.Errorf("failed to start embedded postgres: %w", err)
	}

	cs := fmt.Sprintf("postgres://postgres:postgres@localhost:%d/%s?sslmode=disable", port, dbname)

	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("error creating iofs source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, cs)
	if err != nil {
		return fmt.Errorf("error creating db migrate instance: %w", err)
	}
	m.Log = &migrateLogger{
		logger: logger.Get(ctx),
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running db migrations: %w", err)
	}

	_embeddedPostgresDB = db
	os.Setenv("MODUS_DB", cs)

	logger.Info(ctx).Msg("Embedded PostgreSQL database started successfully.")
	return nil
}

func getEmbeddedPostgresBinariesPath() string {
	return filepath.Join(app.ModusHomeDir(), "postgres")
}

func findAvailablePort(startPort, maxPort uint32) (uint32, error) {
	for port := startPort; port <= maxPort; port++ {
		addr := fmt.Sprintf("localhost:%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available db ports between %d and %d", startPort, maxPort)
}

func shutdownEmbeddedPostgresDB(ctx context.Context) {
	logger.Info(ctx).Msg("Shutting down embedded PostgreSQL database server:")

	if err := _embeddedPostgresDB.Stop(); err != nil {
		logger.Error(ctx).Err(err).Msg("Failed to stop embedded PostgreSQL database.")
	}
}

func shutdownPreviousEmbeddedPostgresDB(ctx context.Context, dataDir string) error {

	// does dataDir exist?
	if _, err := os.Stat(dataDir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error checking dataDir: %w", err)
	}

	// does `postmaster.pid` exist?
	pidFile := filepath.Join(dataDir, "postmaster.pid")
	if _, err := os.Stat(pidFile); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	}

	logger.Warn(ctx).Msg("Previous embedded PostgreSQL instance was not shut down cleanly.  Shutting it down now:")

	// first try to shutdown the process cleanly
	pgctl := filepath.Join(getEmbeddedPostgresBinariesPath(), "bin", "pg_ctl")
	if runtime.GOOS == "windows" {
		pgctl += ".exe"
	}
	if _, err := os.Stat(pgctl); err == nil {
		p := exec.Command(pgctl, "stop", "-w", "-D", dataDir)
		p.Stdout = os.Stdout
		p.Stderr = os.Stderr
		if err := p.Run(); err == nil {
			return nil
		} else {
			logger.Err(ctx, err).Msg("Failed to stop embedded PostgreSQL instance cleanly. Killing db process.")
		}
	}

	// read the pid from the first line of the file
	file, err := os.Open(pidFile)
	if err != nil {
		return fmt.Errorf("error opening pid file: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		pidStr := scanner.Text()
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return fmt.Errorf("error parsing pid: %w", err)
		}

		// kill the process
		process, err := os.FindProcess(pid)
		if err != nil {
			return fmt.Errorf("error finding process: %w", err)
		}
		if err := process.Kill(); err != nil {
			return fmt.Errorf("error killing process: %w", err)
		}
	}

	return nil
}

var migrateLoggerColor = color.New(color.FgWhite, color.Faint)

type migrateLogger struct {
	logger  *zerolog.Logger
	started bool
}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	if !l.started {
		l.logger.Info().Msg("Applying db migrations:")
		l.started = true
	}
	migrateLoggerColor.Fprintf(os.Stderr, "  "+format, v...)
}

func (l *migrateLogger) Verbose() bool {
	return false
}
