package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestEnsureDirCreatesPath(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "deep", "nested", "envm.db")

	if err := ensureDir(dbPath); err != nil {
		t.Fatalf("ensureDir returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Dir(dbPath)); err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}
}

func TestNewSQLiteCreatesDatabaseAndMigrates(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "data", "envm.db")

	db, err := NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("NewSQLite returned error: %v", err)
	}
	defer db.Close()

	if !hasTable(t, db, "projects") {
		t.Fatal("expected project table to be migrated")
	}

	if !hasTable(t, db, "env_vars") {
		t.Fatal("expected env_vars table to be migrated")
	}

	if !hasTable(t, db, "tokens") {
		t.Fatal("expected tokens table to be migrated")
	}

	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected sqlite db file to exist: %v", err)
	}
}

func TestNewSQLiteInvalidPath(t *testing.T) {
	_, err := NewSQLite(t.TempDir())
	if err == nil {
		t.Fatal("expected error for invalid database path")
	}
}

func hasTable(t *testing.T, db *sqlx.DB, tableName string) bool {
	t.Helper()

	var name string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name = ?`, tableName).Scan(&name)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		t.Fatalf("failed checking table existence for %s: %v", tableName, err)
	}

	return name == tableName
}
