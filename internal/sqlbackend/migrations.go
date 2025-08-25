package sqlbackend

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// runMigrations applies all pending migrations to the database
func runMigrations(db *sql.DB) error {
	// Create the source driver from embedded filesystem
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source driver: %w", err)
	}

	// Create the database driver
	dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{
		NoTxWrap: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create the migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", dbDriver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Apply all migrations
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
