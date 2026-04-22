package database

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

const migrationsDir = "migrations"

//go:embed migrations/*.sql
var migrationFiles embed.FS

func NewSQLite(path string) (*sqlx.DB, error) {
	if err := ensureDir(path); err != nil {
		return nil, err
	}

	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	goose.SetBaseFS(migrationFiles)
	if err := goose.SetDialect("sqlite3"); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := goose.Up(db.DB, migrationsDir); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func ensureDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0o700)
}
