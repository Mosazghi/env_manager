package config

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "envmgr"
	keyName     = "master_key"
)

func GetOrCreateMasterKey() ([]byte, error) {
	key, err := getMasterKeyFromKeyring()
	if err == nil {
		return key, nil
	}

	if !errors.Is(err, keyring.ErrNotFound) {
		fallbackKey, fallbackErr := getMasterKeyFromFile()
		if fallbackErr == nil {
			return fallbackKey, nil
		}

		if !errors.Is(fallbackErr, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to read master key from keyring (%v) and fallback file: %w", err, fallbackErr)
		}

		generatedKey, genErr := generateMasterKey()
		if genErr != nil {
			return nil, genErr
		}

		if setErr := setMasterKeyToFile(generatedKey); setErr != nil {
			return nil, fmt.Errorf("keyring unavailable (%v) and failed to write fallback key file: %w", err, setErr)
		}

		return generatedKey, nil
	}

	fallbackKey, fallbackErr := getMasterKeyFromFile()
	if fallbackErr == nil {
		return fallbackKey, nil
	}

	if !errors.Is(fallbackErr, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to read master key fallback file: %w", fallbackErr)
	}

	generatedKey, genErr := generateMasterKey()
	if genErr != nil {
		return nil, genErr
	}

	encoded := base64.StdEncoding.EncodeToString(generatedKey)
	if err := keyring.Set(serviceName, keyName, encoded); err != nil {
		if setErr := setMasterKeyToFile(generatedKey); setErr != nil {
			return nil, fmt.Errorf("failed to store master key in keyring (%v) and fallback file: %w", err, setErr)
		}
	}

	return generatedKey, nil
}

func DeleteMasterKey() error {
	var errs []error

	if err := keyring.Delete(serviceName, keyName); err != nil && !errors.Is(err, keyring.ErrNotFound) {
		errs = append(errs, err)
	}

	if err := deleteMasterKeyFile(); err != nil && !errors.Is(err, os.ErrNotExist) {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func getMasterKeyFromKeyring() ([]byte, error) {
	str, err := keyring.Get(serviceName, keyName)

	if err != nil {
		return nil, err
	}

	key, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("invalid keyring master key encoding: %w", err)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("invalid keyring master key length: %d", len(key))
	}

	return key, nil
}

func generateMasterKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return key, nil
}

func getMasterKeyFromFile() ([]byte, error) {
	path, err := masterKeyFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(data)))
	if err != nil {
		return nil, fmt.Errorf("invalid fallback master key encoding: %w", err)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("invalid fallback master key length: %d", len(key))
	}

	return key, nil
}

func setMasterKeyToFile(key []byte) error {
	path, err := masterKeyFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(key)
	return os.WriteFile(path, []byte(encoded), 0o600)
}

func deleteMasterKeyFile() error {
	path, err := masterKeyFilePath()
	if err != nil {
		return err
	}

	return os.Remove(path)
}

func masterKeyFilePath() (string, error) {
	if path := os.Getenv("ENVM_MASTER_KEY_FILE"); path != "" {
		return path, nil
	}

	configDir, err := os.UserConfigDir()
	if err == nil {
		return filepath.Join(configDir, "envm", "master.key"), nil
	}

	if runtime.GOOS == "linux" {
		return filepath.Join("/var/lib", "envm", "master.key"), nil
	}

	return "", fmt.Errorf("cannot determine master key path: %w", err)
}
