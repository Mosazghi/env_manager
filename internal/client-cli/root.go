package clientcli

import (
	"fmt"
	"log"
	"os"

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

func loadLocalClientConfig(path string) error {
	pID, err := GetStoredProjectID()
	if err != nil {
		return err
	}
	projectID = pID
	return nil
}

func loadGlobalClientConfig() error {
	t, err := GetStoredToken()
	if err != nil {
		return err
	}
	token = t

	u, err := GetStoredServerURL()
	if err != nil {
		return err
	}
	serverURL = u

	return nil
}

func init() {
	if err := loadLocalClientConfig(".envm.config"); err != nil && !os.IsNotExist(err) {
		log.Printf("warning: failed to load .envm.config: %v", err)
	}

	if err := loadGlobalClientConfig(); err != nil && !os.IsNotExist(err) {
		log.Printf("warning: failed to load .envm.config: %v", err)
	}
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", token, "API token")
	rootCmd.PersistentFlags().StringVarP(&projectID, "project-id", "i", projectID, "Default Project ID")
	rootCmd.PersistentFlags().StringVarP(&serverURL, "server-url", "u", serverURL, "Default Server URL")
	rootCmd.PersistentFlags().BoolVarP(&silentMode, "silent-mode", "s", false, "Silent Mode")
}
