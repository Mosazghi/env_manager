package servercli

import (
	"errors"
	"net/http"
	"path/filepath"
	"testing"

	"env-manager/internal/database"
	"env-manager/internal/models"
)

func TestProgramStartAndStop(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "envm.db")
	t.Setenv("PORT", "0")
	t.Setenv("DB_PATH", dbPath)
	t.Setenv("APP_ENV", "test")
	t.Setenv("ENVM_MASTER_KEY_FILE", filepath.Join(t.TempDir(), "master.key"))

	p := &program{}
	if err := p.Start(nil); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}

	if p.srv == nil {
		t.Fatal("expected server to be initialized")
	}

	err := p.Stop(nil)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		t.Fatalf("unexpected stop error: %v", err)
	}
}

func TestTokenCreateCommandRunE(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "envm.db")
	t.Setenv("DB_PATH", dbPath)
	t.Setenv("PORT", "8081")
	t.Setenv("APP_ENV", "test")

	if err := tokenCreateCmd.Flags().Set("expires-in", "5m"); err != nil {
		t.Fatalf("failed to set expires-in flag: %v", err)
	}

	if err := tokenCreateCmd.RunE(tokenCreateCmd, nil); err != nil {
		t.Fatalf("tokenCreateCmd RunE returned error: %v", err)
	}

	db, err := database.NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("failed to reopen sqlite db: %v", err)
	}

	var count int64
	if err := db.Model(&models.Token{}).Count(&count).Error; err != nil {
		t.Fatalf("failed counting tokens: %v", err)
	}

	if count == 0 {
		t.Fatal("expected at least one token created")
	}
}

func TestTokenCreateCommandInvalidDuration(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "envm.db")
	t.Setenv("DB_PATH", dbPath)
	t.Setenv("PORT", "8082")
	t.Setenv("APP_ENV", "test")

	if err := tokenCreateCmd.Flags().Set("expires-in", "invalid"); err != nil {
		t.Fatalf("failed to set expires-in flag: %v", err)
	}

	err := tokenCreateCmd.RunE(tokenCreateCmd, nil)
	if err == nil {
		t.Fatal("expected invalid duration to return error")
	}
}
