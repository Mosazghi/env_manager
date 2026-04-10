package handler

import (
	"env-manager/internal/config"
	"env-manager/internal/crypto"
	"env-manager/internal/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ToResponse(sucess bool, msg string, data any) gin.H {
	return gin.H{"sucess": sucess, "message": msg, "data": data}
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
