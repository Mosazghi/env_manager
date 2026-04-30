package clientcli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
	Use:   "create",
	Short: "Create a new environment variable",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, err := rootCmd.Flags().GetString("server-url")
		if err != nil {
			return fmt.Errorf("failed to get server URL: %w", err)
		}
		client := api.NewClient(token, baseURL)

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			return fmt.Errorf("failed to get key: %w", err)
		}
		value, err := cmd.Flags().GetString("value")
		if err != nil {
			return fmt.Errorf("failed to get value: %w", err)
		}
		pID, err := cmd.Flags().GetString("project-id")
		if err != nil {
			return fmt.Errorf("failed to get project ID: %w", err)
		}
		pIDInt, err := strconv.Atoi(pID)
		if err != nil {
			return fmt.Errorf("project ID conversion failed")
		}

		var body models.CreateEnvVarRequest
		body.ProjectID = pIDInt
		body.Key = key
		body.Value = value

		_, err = client.Post("/env-vars", body)
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
		data, err := client.Get("/projects/"+projectID+"/env-vars", nil)
		if err != nil {
			return err
		}

		var resp struct {
			Data []models.EnvVar `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		localEnvVars, _ := getLocalEnvVars(".env")

		f, err := os.OpenFile(".env", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer f.Close()

		for _, env := range resp.Data {
			if _, exists := localEnvVars[env.Key]; !exists {
				if _, err := fmt.Fprintf(f, "%s=%s\n", env.Key, env.Value); err != nil {
					fmt.Printf("warning: failed to write key '%v' to file", env.Key)
				}
			}
		}

		return nil
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Fetch env vars and run a subcommand",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, _ := rootCmd.Flags().GetString("server-url")
		client := api.NewClient(token, baseURL)
		projectID, _ := rootCmd.Flags().GetString("project-id")
		data, err := client.Get("/projects/"+projectID+"/env-vars", nil)
		if err != nil {
			return err
		}

		var resp struct {
			Data []models.EnvVar `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}
		dashIndex := cmd.ArgsLenAtDash()
		if dashIndex == -1 || dashIndex >= len(args) {
			return fmt.Errorf("no command provided after --")
		}

		commandToRun := args[dashIndex:]
		executable := commandToRun[0]
		remainingArgs := commandToRun[1:]
		execCmd := exec.Command(executable, remainingArgs...)
		newEnv := os.Environ()

		for _, env := range resp.Data {
			newEnv = append(newEnv, fmt.Sprintf("%v=%v", env.Key, env.Value))
		}
		execCmd.Env = newEnv
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		fmt.Println("running cmd: ", commandToRun)
		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("error running %v: %w", executable, err)
		}

		return nil
	},
}

var syncEnvVarsCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync env variables for project",
	RunE: func(clientcli *cobra.Command, args []string) error {
		baseURL, _ := rootCmd.Flags().GetString("server-url")
		client := api.NewClient(token, baseURL)

		filePath, _ := clientcli.Flags().GetString("file-path")
		forceUpdate, err := clientcli.Flags().GetBool("force-update")
		if err != nil {
			return err
		}

		localEnvVars, _ := getLocalEnvVars(filePath)

		projectID, _ := rootCmd.Flags().GetString("project-id")
		silentMode, _ := rootCmd.Flags().GetBool("silent-mode")
		data, err := client.Get("/projects/"+projectID+"/env-vars", nil)
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
				fmt.Printf("%v uploaded\n", key)
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
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("failed to open file %s", err)
		}
		defer f.Close()

		forcePull, err := clientcli.Flags().GetBool("force-pull")
		if err != nil {
			return err
		}
		forcePulledCount := 0
		for key, pair := range remoteEnvVars {
			if _, exists := localEnvVars[key]; !exists {
				if forcePull {
					if _, err := fmt.Fprintf(f, "%s=%s\n", key, pair.Value); err != nil {
						return fmt.Errorf("failed to pull env var: %s", err)
					}
					forcePulledCount++
					continue
				}

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

					if _, err := fmt.Fprintf(f, "%s=%s\n", key, pair.Value); err != nil {
						return fmt.Errorf("failed to pull env var: %s", err)
					}
				}

			}
		}

		if forcePull {
			fmt.Printf("pulled %v env variables\n", forcePulledCount)
		}

		return nil
	},
}

func generateStars(str string) string {
	return strings.Repeat("*", len(str))
}

func init() {
	syncEnvVarsCmd.Flags().Bool("force-update", false, "force variable updates to server")
	syncEnvVarsCmd.Flags().Bool("force-pull", false, "force variable pull from server")
	syncEnvVarsCmd.Flags().StringP("file-path", "p", ".env", "filepath to .env")
	createEnvVarCmd.Flags().StringP("key", "k", "", "env var key")
	createEnvVarCmd.Flags().StringP("value", "v", "", "env var value")
	_ = createEnvVarCmd.MarkFlagRequired("key")
	_ = createEnvVarCmd.MarkFlagRequired("value")
	envVarsCmd.AddCommand(createEnvVarCmd, loadEnvsForProjectCmd, syncEnvVarsCmd, runCmd)
	rootCmd.AddCommand(envVarsCmd)
}
