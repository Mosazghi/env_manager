package config

import (
	"bytes"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetEnv(t *testing.T) {
	t.Setenv("ENVM_TEST_ENV", "")
	if got := getEnv("ENVM_TEST_ENV", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}

	t.Setenv("ENVM_TEST_ENV", "value")
	if got := getEnv("ENVM_TEST_ENV", "fallback"); got != "value" {
		t.Fatalf("expected env value, got %q", got)
	}
}

func TestDefaultDBPath(t *testing.T) {
	if runtime.GOOS == "linux" {
		want := filepath.Join("/var/lib", "envm", "envm.db")
		if got := defaultDBPath(); got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
		return
	}
}

func TestMasterKeyFilePathPriority(t *testing.T) {
	envFile := filepath.Join(t.TempDir(), "custom-master.key")

	t.Setenv("ENVM_MASTER_KEY_FILE", envFile)

	got, err := masterKeyFilePath()
	if err != nil {
		t.Fatalf("masterKeyFilePath returned error: %v", err)
	}

	if got != envFile {
		t.Fatalf("expected env override path %q, got %q", envFile, got)
	}

	t.Setenv("ENVM_MASTER_KEY_FILE", "")
	got, err = masterKeyFilePath()
	if err != nil {
		t.Fatalf("masterKeyFilePath returned error: %v", err)
	}
}

func TestSetGetAndDeleteMasterKeyFile(t *testing.T) {
	keyPath := filepath.Join(t.TempDir(), "nested", "master.key")
	t.Setenv("ENVM_MASTER_KEY_FILE", keyPath)

	key, err := generateMasterKey()
	if err != nil {
		t.Fatalf("generateMasterKey returned error: %v", err)
	}

	if err := setMasterKeyToFile(key); err != nil {
		t.Fatalf("setMasterKeyToFile returned error: %v", err)
	}

	readKey, err := getMasterKeyFromFile()
	if err != nil {
		t.Fatalf("getMasterKeyFromFile returned error: %v", err)
	}

	if !bytes.Equal(readKey, key) {
		t.Fatal("read key does not match written key")
	}

	if err := deleteMasterKeyFile(); err != nil {
		t.Fatalf("deleteMasterKeyFile returned error: %v", err)
	}

	_, err = os.Stat(keyPath)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected file to be removed, got stat error: %v", err)
	}
}

func TestGetOrCreateMasterKeyReturnsKey(t *testing.T) {
	t.Setenv("ENVM_MASTER_KEY_FILE", filepath.Join(t.TempDir(), "master.key"))

	key, err := GetOrCreateMasterKey()
	if err != nil {
		t.Fatalf("GetOrCreateMasterKey returned error: %v", err)
	}

	if len(key) != 32 {
		t.Fatalf("expected 32-byte key, got %d", len(key))
	}
}

func TestLoadUsesEnvironmentValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DB_PATH", "/tmp/custom.db")
	t.Setenv("APP_ENV", "test")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Fatalf("expected port 9090, got %q", cfg.Port)
	}

	if cfg.DBPath != "/tmp/custom.db" {
		t.Fatalf("expected db path /tmp/custom.db, got %q", cfg.DBPath)
	}

	if cfg.Env != "test" {
		t.Fatalf("expected env test, got %q", cfg.Env)
	}
}

func TestGetOrCreateMasterKeyUsesExistingFallbackFile(t *testing.T) {
	keyPath := filepath.Join(t.TempDir(), "master.key")
	t.Setenv("ENVM_MASTER_KEY_FILE", keyPath)

	original := bytes.Repeat([]byte{1}, 32)
	encoded := base64.StdEncoding.EncodeToString(original)
	if err := os.WriteFile(keyPath, []byte(encoded), 0o600); err != nil {
		t.Fatalf("failed to write key file: %v", err)
	}

	key, err := GetOrCreateMasterKey()
	if err != nil {
		t.Fatalf("GetOrCreateMasterKey returned error: %v", err)
	}

	if len(key) != 32 {
		t.Fatalf("expected 32-byte key, got %d", len(key))
	}
}

func TestDeleteMasterKeyRemovesFallbackFile(t *testing.T) {
	keyPath := filepath.Join(t.TempDir(), "master.key")
	t.Setenv("ENVM_MASTER_KEY_FILE", keyPath)

	if err := os.WriteFile(keyPath, []byte(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{2}, 32))), 0o600); err != nil {
		t.Fatalf("failed to write key file: %v", err)
	}

	_ = DeleteMasterKey()

	_, err := os.Stat(keyPath)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected key file to be removed, got stat error: %v", err)
	}
}
