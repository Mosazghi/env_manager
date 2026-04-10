package clientcli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"env-manager/internal/api"
	"env-manager/internal/models"

	"github.com/spf13/cobra"
)

var envVarsCmd = &cobra.Command{
	Use:   "env-vars",
	Short: "Manage environment variables",
}

var createEnvVarCmd = &cobra.Command{
	Use:   "create [key] [value]",
	Short: "Create a new environment variable",
	Args:  cobra.ExactArgs(2),
	RunE: func(clientcli *cobra.Command, args []string) error {
		baseURL, _ := rootCmd.Flags().GetString("server-url")
		client := api.NewClient(token, baseURL)
		var body models.CreateEnvVarRequest
		body.Key = args[0]
		body.Value = args[1]
		_, err := client.Post("/env-vars", body)
		if err != nil {
			return err
		}
		fmt.Println("Environment variable created")
		return nil
	},
}

var loadEnvsForProjectCmd = &cobra.Command{
	Use:   "load",
	Short: "Load env variables for project",
	RunE: func(clientcli *cobra.Command, args []string) error {
		baseURL, _ := rootCmd.Flags().GetString("server-url")
		client := api.NewClient(token, baseURL)
		projectID, _ := rootCmd.Flags().GetString("project-id")
		data, err := client.Get("/projects/" + projectID + "/env-vars")
		if err != nil {
			return err
		}

		var resp struct {
			Data []models.EnvVar `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		f, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}

		defer f.Close()

		for _, env := range resp.Data {
			if _, err := fmt.Fprintf(f, "%s=%s\n", env.Key, env.Value); err != nil {
				fmt.Printf("warning: failed to write key '%v' to file", env.Key)
			}
		}

		return nil
	},
}

var syncEnvVarsCmd = &cobra.Command{
	Use:   "sync [force] [filePath]",
	Short: "Sync env variables for project",
	RunE: func(clientcli *cobra.Command, args []string) error {
		baseURL, _ := rootCmd.Flags().GetString("server-url")
		client := api.NewClient(token, baseURL)

		filePath, _ := clientcli.Flags().GetString("file-path")
		forceUpdate, err := clientcli.Flags().GetBool("force-update")
		if err != nil {
			return err
		}

		localEnvVars := make(map[string]string)

		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0o644)
		if err != nil {
			return fmt.Errorf("failed to open file: %s", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			variable := strings.Split(line, "=")
			key := variable[0]
			val := variable[1]
			if key != "" && val != "" {
				localEnvVars[key] = val
			}
		}

		projectID, _ := rootCmd.Flags().GetString("project-id")
		silentMode, _ := rootCmd.Flags().GetBool("silent-mode")
		data, err := client.Get("/projects/" + projectID + "/env-vars")
		if err != nil {
			return err
		}

		var resp struct {
			Data []models.EnvVar `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		type IDValPair struct {
			ID    uint
			Value string
		}
		remoteEnvVars := make(map[string]IDValPair)
		for _, env := range resp.Data {
			remoteEnvVars[env.Key] = IDValPair{ID: env.ID, Value: env.Value}
		}

		for key, val := range localEnvVars {
			// If it doesn't remote, create it
			remoteEnvVar, exist := remoteEnvVars[key]
			if !exist {
				pIDInt, err := strconv.Atoi(projectID)
				if err != nil {
					return fmt.Errorf("project ID conversion failed")
				}
				if _, err := client.Post("/env-vars", models.CreateEnvVarRequest{Key: key, Value: val, ProjectID: pIDInt}); err != nil {
					return fmt.Errorf("failed to create env var: %s", err)
				}
				fmt.Println("uploaded a new variable")
			}

			// If it exists but value is different, update it
			if exist {
				if remoteEnvVar.Value != val {
					update := func() error {
						if _, err := client.Put("/env-vars/"+fmt.Sprint(remoteEnvVar.ID), models.UpdateEnvVarRequest{Value: val}); err != nil {
							return fmt.Errorf("failed to update env var: %s", err)
						}
						return nil
					}

					if forceUpdate {
						if err := update(); err != nil {
							return err
						}
					} else {
						var confirmation string

						msg := "%v's value changed: %v => %v. Update to remote (y/N)? "
						if silentMode {
							fmt.Printf(msg, key, generateStars(remoteEnvVar.Value), generateStars(val))
						} else {
							fmt.Printf(msg, key, remoteEnvVar.Value, val)
						}

						fmt.Scanln(&confirmation)

						confirmation = strings.ToLower(confirmation)

						if confirmation == "y" {
							if err := update(); err != nil {
								return err
							}
						}
					}
				}
			}
		}

		// if it exists on remote, but not locally, then ask user if to delete or keep
		for key, pair := range remoteEnvVars {
			if _, exists := localEnvVars[key]; !exists {
				var confirmation string
				msg := "%v=%v doesn't exists locally. Actions: (delete=d, pull=p, nothing=N)?"
				if silentMode {
					fmt.Printf(msg, key, generateStars(pair.Value))
				} else {
					fmt.Printf(msg, key, pair.Value)
				}
				fmt.Scanln(&confirmation)

				confirmation = strings.ToLower(confirmation)

				switch confirmation {

				case "d":
					if _, err := client.Delete("/env-vars/" + fmt.Sprint(pair.ID)); err != nil {
						fmt.Printf("failed to delete env var: %s\n", err)
					}
				case "p":
					f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
					if err != nil {
						return fmt.Errorf("failed to open file %s", err)
					}

					defer f.Close()

					if _, err := fmt.Fprintf(f, "%s=%s\n", key, pair.Value); err != nil {
						return fmt.Errorf("failed to pull env var: %s", err)
					}
				}

			}
		}

		return nil
	},
}

func generateStars(str string) string {
	return strings.Repeat("*", len(str))
}

func init() {
	syncEnvVarsCmd.Flags().BoolP("force-update", "f", false, "force variable updates")
	syncEnvVarsCmd.Flags().StringP("file-path", "p", ".env", "filepath to .env")
	envVarsCmd.AddCommand(createEnvVarCmd, loadEnvsForProjectCmd, syncEnvVarsCmd)
	rootCmd.AddCommand(envVarsCmd)
}
