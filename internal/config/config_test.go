package config

import (
	"bytes"
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

	stateDir := t.TempDir()
	t.Setenv("STATE_DIRECTORY", stateDir)

	want := filepath.Join(stateDir, "envm.db")
	if got := defaultDBPath(); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestMasterKeyFilePathPriority(t *testing.T) {
	stateDir := t.TempDir()
	envFile := filepath.Join(t.TempDir(), "custom-master.key")

	t.Setenv("STATE_DIRECTORY", stateDir)
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

	want := filepath.Join(stateDir, "master.key")
	if got != want {
		t.Fatalf("expected state directory path %q, got %q", want, got)
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

func TestGetMasterKeyFromFileInvalidEncoding(t *testing.T) {
	keyPath := filepath.Join(t.TempDir(), "master.key")
	t.Setenv("ENVM_MASTER_KEY_FILE", keyPath)

	if err := os.WriteFile(keyPath, []byte("not-valid-base64"), 0o600); err != nil {
		t.Fatalf("failed to write invalid key file: %v", err)
	}

	_, err := getMasterKeyFromFile()
	if err == nil {
		t.Fatal("expected error for invalid key encoding")
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
