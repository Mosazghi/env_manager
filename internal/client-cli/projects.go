package clientcli

import (
	"encoding/json"
	"fmt"
	"os"
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
	RunE: func(clientcli *cobra.Command, args []string) error {
		baseURL, _ := rootCmd.Flags().GetString("server-url")
		client := api.NewClient(token, baseURL)
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
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t\n", p.ID, truncateProjectName(p.Name), truncateProjectDescription(p.Description), p.CreatedAt.Format("2006-01-02 15:04:05"), p.UpdatedAt.Format("2006-01-02 15:04:05"))
		}

		w.Flush()
		return nil
	},
}

var createProjectCmd = &cobra.Command{
	Use:   "create [name] [description]",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(2),
	RunE: func(clientcli *cobra.Command, args []string) error {
		baseURL, _ := rootCmd.Flags().GetString("server-url")
		client := api.NewClient(token, baseURL)
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

func init() {
	projectsCmd.AddCommand(fetchProjectsCmd, createProjectCmd)
	rootCmd.AddCommand(projectsCmd)
}
