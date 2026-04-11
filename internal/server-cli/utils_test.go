package servercli

import (
	"os"
	"path/filepath"
	"strings"
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

func TestDeliverTokenFileMode(t *testing.T) {
	tokenPath := filepath.Join(t.TempDir(), "token.txt")

	if err := deliverToken("abc123", "1h", "file", tokenPath); err != nil {
		t.Fatalf("deliverToken returned error: %v", err)
	}

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatalf("failed reading token file: %v", err)
	}

	if strings.TrimSpace(string(data)) != "abc123" {
		t.Fatalf("unexpected token file content: %q", string(data))
	}
}

func TestDeliverTokenInvalidMode(t *testing.T) {
	err := deliverToken("abc123", "1h", "invalid", "")
	if err == nil {
		t.Fatal("expected error for invalid output mode")
	}
}
