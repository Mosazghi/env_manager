package handler

import (
	"encoding/json"
	"env-manager/internal/config"
	"env-manager/internal/crypto"
	"env-manager/internal/models"
	"fmt"
	"net/http"
)

func ToResponse(sucess bool, msg string, data any) map[string]any {
	return map[string]any{"sucess": sucess, "message": msg, "data": data}
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "{\"error\":\"failed to encode response\"}", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func DecodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is required")
	}

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return err
	}

	return nil
}

func DecryptEnvVars(envVars *[]models.EnvVar) ([]models.EnvVar, error) {
	output := make([]models.EnvVar, len(*envVars))

	masterKey, err := config.GetOrCreateMasterKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get master key: %v", err)
	}
	copy(output, *envVars)

	for i := range output {
		dec, err := crypto.Decrypt(masterKey, output[i].EncryptedVal)
		if err != nil {
			fmt.Printf("error decrypting %v: %v", output[i].Key, err)
			continue
		}

		output[i].Value = string(dec)
	}

	return output, nil
}

func EncryptEnvVars(envVars *[]models.EnvVar) ([]models.EnvVar, error) {
	output := make([]models.EnvVar, len(*envVars))

	masterKey, err := config.GetOrCreateMasterKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get master key: %v", err)
	}
	copy(output, *envVars)

	for i := range output {
		enc, err := crypto.Encrypt(masterKey, []byte(output[i].EncryptedVal))
		if err != nil {
			fmt.Printf("error encrypting %v: %v", output[i].Key, err)
			continue
		}

		output[i].EncryptedVal = enc

	}

	return output, nil
}

func EncryptValue(val string) (string, error) {
	masterKey, err := config.GetOrCreateMasterKey()
	if err != nil {
		return "", fmt.Errorf("failed to get master key: %v", err)
	}

	return crypto.Encrypt(masterKey, []byte(val))
}

func DecryptValue(val string) ([]byte, error) {
	masterKey, err := config.GetOrCreateMasterKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get master key: %v", err)
	}

	return crypto.Decrypt(masterKey, val)
}
