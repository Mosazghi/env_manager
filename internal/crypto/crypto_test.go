package crypto

import (
	"bytes"
	"testing"
)

var testKey = []byte("0123456789abcdef0123456789abcdef")

func TestEncryptDecryptRoundTrip(t *testing.T) {
	plain := []byte("postgres://localhost")

	encrypted, err := Encrypt(testKey, plain)
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}

	decrypted, err := Decrypt(testKey, encrypted)
	if err != nil {
		t.Fatalf("Decrypt returned error: %v", err)
	}

	if !bytes.Equal(decrypted, plain) {
		t.Fatalf("expected %q, got %q", string(plain), string(decrypted))
	}
}

func TestEncryptInvalidKeyLength(t *testing.T) {
	_, err := Encrypt([]byte("short"), []byte("value"))
	if err == nil {
		t.Fatal("expected error for invalid key length")
	}
}

func TestDecryptInvalidBase64(t *testing.T) {
	_, err := Decrypt(testKey, "not-base64$$")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecryptInvalidCiphertextLength(t *testing.T) {
	_, err := Decrypt(testKey, "YWJj")
	if err == nil {
		t.Fatal("expected error for invalid ciphertext")
	}
}

func TestDecryptWithWrongKeyFails(t *testing.T) {
	plain := []byte("top-secret")
	encrypted, err := Encrypt(testKey, plain)
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}

	wrongKey := []byte("fedcba9876543210fedcba9876543210")
	_, err = Decrypt(wrongKey, encrypted)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptInvalidKeyLength(t *testing.T) {
	enc, err := Encrypt(testKey, []byte("value"))
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}

	_, err = Decrypt([]byte("short"), enc)
	if err == nil {
		t.Fatal("expected decrypt to fail for invalid key length")
	}
}
