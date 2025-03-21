// Package migrator is responsible to run migrations.
package migrator

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Migrator structure of object that holds connection to DB.
type Migrator struct {
	db *sql.DB
}

// NewMigrator function for Migrator construction.
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db: db,
	}
}

// RunMigrations function launches migrations for DB
func (m *Migrator) RunMigrations() error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(m.db, "migrations"); err != nil {
		return err
	}

	return nil
}
