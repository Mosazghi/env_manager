package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	. "env-manager/internal/shared"

	"github.com/spf13/cobra"
)

var token = "4ebe17469d06d6823d9e9339ae97085d2c8bbca82f5e559ac3a48b6ecd7e8e67c20c2f35ef62c313e7eb752f42ff9525"

func main() {
	cmdCreateProject := &cobra.Command{
		Use:   "create-project",
		Short: "Create new project(s)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range args {
				log.Printf("Creating for %v", name)
				req, err := http.NewRequest("POST", "http://localhost:8080/env/projects/"+name, nil)
				if err != nil {
					log.Fatal(err)
				}

				req.Header.Add("Authorization", "Bearer "+token)

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatalf("Error reading body: %v", err)
				}

				fmt.Printf("Response: %v\n", string(body))
			}
		},
	}

	cmdFetchProjects := &cobra.Command{
		Use:   "fetch-projects",
		Short: "Fetch all projects from the server",
		Long:  "Fetch all projects from the server",
		Run: func(cmd *cobra.Command, args []string) {
			req, err := http.NewRequest("GET", "http://localhost:8080/env/projects", nil)
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Add("Authorization", "Bearer "+token)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			var projects []Project
			err = json.NewDecoder(resp.Body).Decode(&projects)
			if err != nil {
				log.Fatal(err)
			}

			prettyJSON, err := json.MarshalIndent(projects, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Projects: %v", string(prettyJSON))
			err = os.WriteFile("test.json", prettyJSON, 0o644)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdFetchProjects, cmdCreateProject)
	rootCmd.Execute()
}
