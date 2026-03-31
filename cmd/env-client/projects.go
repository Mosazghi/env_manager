package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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
			Data []map[string]any `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		out, _ := json.MarshalIndent(resp.Data, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

var createProjectCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(token)
		data, err := client.Post("/projects/"+args[0], nil)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

var getEnvVarsForProjectCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			log.Fatal("Project ID is required")
		}

		client := api.NewClient(token)
		projectID := args[0]

		data, err := client.Get(fmt.Sprintf("http://localhost:8080/api/projects/%s/env-vars", projectID))
		if err != nil {
			log.Fatal(err)
		}

		var resp struct {
			Data []models.EnvVar `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		envVars := resp.Data

		// write to .env file
		f, err := os.Create(".env")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		for _, envVar := range envVars {
			_, err := f.WriteString(fmt.Sprintf("%s=%s\n", envVar.Key, envVar.Value))
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	},
}

func init() {
	projectsCmd.AddCommand(fetchProjectsCmd, createProjectCmd)
	rootCmd.AddCommand(projectsCmd)
}
