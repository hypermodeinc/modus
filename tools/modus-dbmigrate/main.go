package main

import (
	"embed"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	connStr := os.Getenv("MODUS_DB")
	if connStr == "" {
		fmt.Fprintln(os.Stderr, "MODUS_DB environment variable is not set.")
		os.Exit(1)
	}

	if err := migrateDB(connStr); errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("Database schema is already up to date.")
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Error running migrations: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Migrations applied successfully")
	}
}

func migrateDB(connStr string) error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("error creating iofs source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, connStr)
	if err != nil {
		return err
	}
	m.Log = &migrateLogger{}

	return m.Up()
}

type migrateLogger struct {
}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *migrateLogger) Verbose() bool {
	return false
}
