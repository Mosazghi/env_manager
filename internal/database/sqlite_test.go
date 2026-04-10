package database

import (
	"os"
	"path/filepath"
	"testing"

	"env-manager/internal/models"
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

	if !db.Migrator().HasTable(&models.Project{}) {
		t.Fatal("expected project table to be migrated")
	}

	if !db.Migrator().HasTable(&models.EnvVar{}) {
		t.Fatal("expected env_vars table to be migrated")
	}

	if !db.Migrator().HasTable(&models.Token{}) {
		t.Fatal("expected tokens table to be migrated")
	}

	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected sqlite db file to exist: %v", err)
	}
}
