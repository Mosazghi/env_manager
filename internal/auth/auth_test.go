package auth

import (
	"encoding/hex"
	"testing"
)

func TestGenerateAuthTokenFormat(t *testing.T) {
	token, err := GenerateAuthToken("test-passphrase")
	if err != nil {
		t.Fatalf("GenerateAuthToken returned error: %v", err)
	}

	if len(token) != 96 {
		t.Fatalf("expected token length 96, got %d", len(token))
	}

	raw, err := hex.DecodeString(token)
	if err != nil {
		t.Fatalf("token should be valid hex: %v", err)
	}

	if len(raw) != 48 {
		t.Fatalf("expected decoded token length 48, got %d", len(raw))
	}
}

func TestGenerateAuthTokenIsRandom(t *testing.T) {
	tokenA, err := GenerateAuthToken("test-passphrase")
	if err != nil {
		t.Fatalf("GenerateAuthToken returned error: %v", err)
	}

	tokenB, err := GenerateAuthToken("test-passphrase")
	if err != nil {
		t.Fatalf("GenerateAuthToken returned error: %v", err)
	}

	if tokenA == tokenB {
		t.Fatal("expected different tokens for two calls with random salt")
	}
}
