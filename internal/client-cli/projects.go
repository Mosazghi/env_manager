package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"env-manager/internal/api"
	"env-manager/internal/models"

	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
}

var fetchProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(token)
		data, err := client.Get("/projects")
		if err != nil {
			return err
		}

		var resp struct {
			Data []models.Project `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tDescription\tCreatedAt\tUpdatedAt\t")
		for _, p := range resp.Data {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t\n", p.ID, p.Name, p.Description, p.CreatedAt.Format("2006-01-02 15:04:05"), p.UpdatedAt.Format("2006-01-02 15:04:05"))
		}

		w.Flush()
		return nil
	},
}

var createProjectCmd = &cobra.Command{
	Use:   "create [name] [description]",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(token)
		var body models.CreateProjectRequest
		body.Name = args[0]
		body.Description = args[1]
		_, err := client.Post("/projects/", body)
		if err != nil {
			return err
		}
		fmt.Println("Project created")
		return nil
	},
}

var loadEnvsForProjectCmd = &cobra.Command{
	Use:   "load",
	Short: "Load env variables for project",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(token)
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
			os.Exit(1)
		}

		defer f.Close()

		for _, env := range resp.Data {
			if _, err := fmt.Fprintf(f, "%s=%s\n", env.Key, env.Value); err != nil {
				os.Exit(1)
			}
		}

		return nil
	},
}

var syncEnvVarsCmd = &cobra.Command{
	Use:   "sync [force] [filePath]",
	Short: "Sync env variables for project",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(token)

		filePath, err := cmd.Flags().GetString("file-path")
		forceUpdate, err := cmd.Flags().GetBool("force-update")
		if err != nil {
			return err
		}

		localEnvVars := make(map[string]string)

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("failed to open file: %s\n", err)
			os.Exit(1)
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

		// Remote vars
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
					fmt.Println("Project ID conversion failed")
					os.Exit(1)
				}
				if _, err := client.Post("/env-vars", models.CreateEnvVarRequest{Key: key, Value: val, ProjectID: pIDInt}); err != nil {
					fmt.Printf("failed to create env var: %s\n", err)
					os.Exit(1)
				}
				fmt.Println("Uploaded a new variable")
			}

			// If it exists but value is different, update it
			if exist {
				if remoteEnvVar.Value != val {
					update := func() {
						if _, err := client.Put("/env-vars/"+fmt.Sprint(remoteEnvVar.ID), models.UpdateEnvVarRequest{Value: val}); err != nil {
							fmt.Printf("failed to update env var: %s\n", err)
							os.Exit(1)
						}
					}

					if forceUpdate {
						update()
					} else {
						var confirmation string
						if !silentMode {
							fmt.Printf("%v's value changed: %v => %v. Update to remote (y/N)? ", key, remoteEnvVar.Value, val)
						} else {
							fmt.Printf("%v's value changed: [len=%v] => [len=%v]. Update to remote (y/N)? ", key, len(remoteEnvVar.Value), len(val))
						}
						fmt.Scanln(&confirmation)

						confirmation = strings.ToLower(confirmation)

						if confirmation == "y" {
							update()
						}
					}
				}
			}
		}

		// if it exists on remote, but not locally, then ask user if to delete or keep
		for key, pair := range remoteEnvVars {
			if _, exists := localEnvVars[key]; !exists {
				var confirmation string
				if !silentMode {
					fmt.Printf("%v=%v (%v) doesn't exists locally. Actions: (delete=d, pull=p, nothing=N)? ", key, pair.Value, pair.ID)
				} else {
					fmt.Printf("%v = [len=%v] (%v) doesn't exists locally. Actions: (delete=d, pull=p, nothing=N)? ", key, len(pair.Value), pair.ID)
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
						fmt.Printf("failed to open file %s\n", err)
						os.Exit(1)
					}

					defer f.Close()

					if _, err := fmt.Fprintf(f, "%s=%s\n", key, pair.Value); err != nil {
						fmt.Printf("failed to pull env var: %s\n", err)
						os.Exit(1)
					}
				}

			}
		}

		return nil
	},
}

func init() {
	syncEnvVarsCmd.Flags().Bool("force-update", false, "force variable updates")
	syncEnvVarsCmd.Flags().String("file-path", ".env", "filepath to .env")
	projectsCmd.AddCommand(fetchProjectsCmd, createProjectCmd, loadEnvsForProjectCmd, syncEnvVarsCmd)
	rootCmd.AddCommand(projectsCmd)
}
