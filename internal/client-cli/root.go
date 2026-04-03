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
	Use:   "envm-client",
	Short: "env-manager Client CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	file, err := os.Open(".envm.config")
	if err != nil {
		log.Fatalf("impossible to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "ENVM_TOKEN") {
			token = strings.Split(line, "=")[1]
		}

		if strings.Contains(line, "PROJET_ID") {
			projectID = strings.Split(line, "=")[1]
		}

		if strings.Contains(line, "SERVER_URL") {
			serverURL = strings.Split(line, "=")[1]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("scanner encountered an error: %s", err)
	}

	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", token, "API token (or set ENVM_TOKEN)")
	rootCmd.PersistentFlags().StringVarP(&projectID, "project-id", "i", projectID, "Default Project ID")
	rootCmd.PersistentFlags().StringVarP(&serverURL, "server-url", "u", serverURL, "Default Server Url")
	rootCmd.PersistentFlags().BoolVarP(&silentMode, "silent-mode", "s", false, "Silent Mode")
}
