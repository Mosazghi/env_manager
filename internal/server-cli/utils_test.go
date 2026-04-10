package servercli

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	got, err := parseDuration("2d")
	if err != nil {
		t.Fatalf("parseDuration returned error: %v", err)
	}

	if got != 48*time.Hour {
		t.Fatalf("expected 48h, got %s", got)
	}

	got, err = parseDuration("30m")
	if err != nil {
		t.Fatalf("parseDuration returned error: %v", err)
	}

	if got != 30*time.Minute {
		t.Fatalf("expected 30m, got %s", got)
	}
}

func TestParseDurationInvalid(t *testing.T) {
	_, err := parseDuration("abc")
	if err == nil {
		t.Fatal("expected parseDuration to fail for invalid input")
	}
}

func TestGenerateRandomToken(t *testing.T) {
	a := generateRandomToken()
	b := generateRandomToken()

	if a == "" || b == "" {
		t.Fatal("expected generated token to be non-empty")
	}

	if a == b {
		t.Fatal("expected two generated tokens to differ")
	}
}

func TestGetMasterPassphrase(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CREDENTIALS_DIRECTORY", dir)

	secretPath := filepath.Join(dir, "envm-passphrase")
	if err := os.WriteFile(secretPath, []byte("passphrase-value\n"), 0o600); err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	got, err := getMasterPassphrase()
	if err != nil {
		t.Fatalf("getMasterPassphrase returned error: %v", err)
	}

	if got != "passphrase-value" {
		t.Fatalf("expected trimmed passphrase, got %q", got)
	}
}
