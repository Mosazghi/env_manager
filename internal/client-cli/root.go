package clientcli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	token      string
	projectID  string
	serverURL  string
	silentMode bool
)

var rootCmd = &cobra.Command{
	Use:           "envm-client",
	Short:         "env-manager Client CLI",
	SilenceErrors: true,
	SilenceUsage:  false,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadClientConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]
		switch key {
		case "ENVM_TOKEN":
			token = value
		case "PROJECT_ID", "PROJET_ID":
			projectID = value
		case "SERVER_URL":
			serverURL = value
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner encountered an error: %w", err)
	}

	return nil
}

func init() {
	if err := loadClientConfig(".envm.config"); err != nil && !os.IsNotExist(err) {
		log.Printf("warning: failed to load .envm.config: %v", err)
	}

	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", token, "API token (or set ENVM_TOKEN)")
	rootCmd.PersistentFlags().StringVarP(&projectID, "project-id", "i", projectID, "Default Project ID")
	rootCmd.PersistentFlags().StringVarP(&serverURL, "server-url", "u", serverURL, "Default Server Url")
	rootCmd.PersistentFlags().BoolVarP(&silentMode, "silent-mode", "s", false, "Silent Mode")
}
