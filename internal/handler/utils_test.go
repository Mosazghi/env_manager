package handler

import (
	"path/filepath"
	"testing"

	"env-manager/internal/models"
)

func TestToResponse(t *testing.T) {
	resp := ToResponse(true, "ok", 123)

	if got := resp["sucess"]; got != true {
		t.Fatalf("expected sucess=true, got %v", got)
	}

	if got := resp["message"]; got != "ok" {
		t.Fatalf("expected message=ok, got %v", got)
	}

	if got := resp["data"]; got != 123 {
		t.Fatalf("expected data=123, got %v", got)
	}
}

func TestEncryptAndDecryptValueRoundTrip(t *testing.T) {
	t.Setenv("ENVM_MASTER_KEY_FILE", filepath.Join(t.TempDir(), "master.key"))

	encrypted, err := EncryptValue("secret-value")
	if err != nil {
		t.Fatalf("EncryptValue returned error: %v", err)
	}

	decrypted, err := DecryptValue(encrypted)
	if err != nil {
		t.Fatalf("DecryptValue returned error: %v", err)
	}

	if string(decrypted) != "secret-value" {
		t.Fatalf("expected decrypted value secret-value, got %s", string(decrypted))
	}
}

func TestEncryptAndDecryptEnvVarsRoundTrip(t *testing.T) {
	t.Setenv("ENVM_MASTER_KEY_FILE", filepath.Join(t.TempDir(), "master.key"))

	input := []models.EnvVar{
		{Key: "A", EncryptedVal: "value-a"},
		{Key: "B", EncryptedVal: "value-b"},
	}

	encrypted, err := EncryptEnvVars(&input)
	if err != nil {
		t.Fatalf("EncryptEnvVars returned error: %v", err)
	}

	if encrypted[0].EncryptedVal == "value-a" {
		t.Fatal("expected encrypted value to differ from plaintext")
	}

	decrypted, err := DecryptEnvVars(&encrypted)
	if err != nil {
		t.Fatalf("DecryptEnvVars returned error: %v", err)
	}

	if decrypted[0].Value != "value-a" || decrypted[1].Value != "value-b" {
		t.Fatalf("unexpected decrypted values: %+v", decrypted)
	}
}

func TestDecryptEnvVarsSkipsInvalidCiphertext(t *testing.T) {
	t.Setenv("ENVM_MASTER_KEY_FILE", filepath.Join(t.TempDir(), "master.key"))

	input := []models.EnvVar{{Key: "BROKEN", EncryptedVal: "%%%not-base64%%%"}}
	decrypted, err := DecryptEnvVars(&input)
	if err != nil {
		t.Fatalf("DecryptEnvVars returned error: %v", err)
	}

	if decrypted[0].Value != "" {
		t.Fatalf("expected empty value for invalid ciphertext, got %q", decrypted[0].Value)
	}
}
