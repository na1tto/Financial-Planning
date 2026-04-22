package db

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(database *sql.DB) error {
	migrationsFS, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Up(database, ".")
}
