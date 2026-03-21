package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateAuthToken(passphrase string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(passphrase))
	h.Write(salt)

	token := append(salt, h.Sum(nil)...)

	return hex.EncodeToString(token), nil
}
