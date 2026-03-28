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
		log.Printf("USING TOKEN: %v", token)
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
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(token)
		var body models.CreateProjectRequest
		body.Name = args[0]
		body.Description = args[1]
		data, err := client.Post("/projects/", body)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

var loadEnvsForProject = &cobra.Command{
	Use:   "load [project-id]",
	Short: "Load env variables for project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(token)
		data, err := client.Get("/projects/" + args[0] + "/env-vars")
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
			log.Fatal(err)
		}

		defer f.Close()

		for _, env := range resp.Data {
			if _, err := fmt.Fprintf(f, "%s=%s\n", env.Key, env.Value); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println(string(data))
		return nil
	},
}

func init() {
	projectsCmd.AddCommand(fetchProjectsCmd, createProjectCmd, loadEnvsForProject)
	rootCmd.AddCommand(projectsCmd)
}
