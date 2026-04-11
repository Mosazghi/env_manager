package clientcli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure and see token",
}

type GlobalClientConfig struct {
	Token     string `json:"token"`
	ServerURL string `json:"serverUrl"`
}

type LocalClientConfig struct {
	ProjectID string `json:"projectID"`
}

var storeTokenCmd = &cobra.Command{
	Use:   "store-token [token]",
	Short: "Store a new token used for auth",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := args[0]
		if err := storeToken(token, getGlobalConfigFilePath()); err != nil {
			return fmt.Errorf("failed to store token: %w", err)
		}
		return nil
	},
}

var storeServerURLCmd = &cobra.Command{
	Use:   "store-url [server-url]",
	Short: "Store server url",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		if err := storeServerURL(url, getGlobalConfigFilePath()); err != nil {
			return fmt.Errorf("failed to store url: %w", err)
		}
		return nil
	},
}

var storeProjectID = &cobra.Command{
	Use:   "set-project-id [project-id]",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	RunE: func(clientcli *cobra.Command, args []string) error {
		pID := args[0]
		if err := setProjectID(pID, getLocalConfigPath()); err != nil {
			return fmt.Errorf("failed to project id: %w", err)
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(storeTokenCmd, storeServerURLCmd, storeProjectID)
	rootCmd.AddCommand(configCmd)
}

func setProjectID(id string, filePath string) error {
	config, err := getLocalConfigData(filePath)
	if err != nil {
		return err
	}
	config.ProjectID = id
	jsonData, _ := json.MarshalIndent(config, "", "  ")
	return os.WriteFile(filePath, jsonData, 0o644)
}

func storeToken(token, filePath string) error {
	config, err := getGlobalConfigData(filePath)
	if err != nil {
		return err
	}

	config.Token = token
	jsonData, _ := json.MarshalIndent(config, "", "  ")
	return os.WriteFile(filePath, jsonData, 0o644)
}

func storeServerURL(url, filePath string) error {
	config, err := getGlobalConfigData(filePath)
	if err != nil {
		return err
	}
	config.ServerURL = url
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0o644)
}

func getGlobalConfigData(filePath string) (*GlobalClientConfig, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config GlobalClientConfig

	content, err := io.ReadAll(file)

	if len(content) <= 0 {
		return &GlobalClientConfig{}, nil
	}

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func getLocalConfigData(filePath string) (*LocalClientConfig, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config LocalClientConfig

	content, err := io.ReadAll(file)

	if len(content) <= 0 {
		return &LocalClientConfig{}, nil
	}

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func GetStoredToken() (string, error) {
	config, err := getGlobalConfigData(getGlobalConfigFilePath())
	if err != nil {
		return "", err
	}
	return config.Token, nil
}

func GetStoredServerURL() (string, error) {
	config, err := getGlobalConfigData(getGlobalConfigFilePath())
	if err != nil {
		return "", err
	}
	return config.ServerURL, nil
}

func GetStoredProjectID() (string, error) {
	config, err := getLocalConfigData(getLocalConfigPath())
	if err != nil {
		return "", err
	}
	return config.ProjectID, nil
}

func getGlobalConfigFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("cannot determine config directory: ", err)
	}

	return filepath.Join(configDir, "envm", "envm.global.json")
}

func getLocalConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("cannot determine config directory: ", err)
	}

	return filepath.Join(cwd, ".envm.local.json")
}
