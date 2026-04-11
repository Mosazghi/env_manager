package clientcli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadClientConfig(t *testing.T) {
	token = ""
	projectID = ""
	serverURL = ""

	configPath := filepath.Join(t.TempDir(), ".envm.config")
	content := "ENVM_TOKEN=test-token\nPROJECT_ID=42\nSERVER_URL=http://localhost:8080\n"
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	if err := loadLocalClientConfig(configPath); err != nil {
		t.Fatalf("loadClientConfig returned error: %v", err)
	}

	if token != "test-token" {
		t.Fatalf("expected token test-token, got %q", token)
	}

	if projectID != "42" {
		t.Fatalf("expected projectID 42, got %q", projectID)
	}

	if serverURL != "http://localhost:8080" {
		t.Fatalf("expected serverURL http://localhost:8080, got %q", serverURL)
	}
}

func TestLoadClientConfigMissingFile(t *testing.T) {
	err := loadLocalClientConfig(filepath.Join(t.TempDir(), "missing.config"))
	if err == nil {
		t.Fatal("expected error when config file is missing")
	}
}

func TestGenerateStars(t *testing.T) {
	if got := generateStars("abcd"); got != "****" {
		t.Fatalf("expected ****, got %q", got)
	}

	if got := generateStars(""); got != "" {
		t.Fatalf("expected empty output for empty input, got %q", got)
	}
}
