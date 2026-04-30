package clientcli

import (
	"encoding/json"
	"env-manager/internal/api"
	"env-manager/internal/models"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var tokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "Manage tokens",
}

var listTokensCmd = &cobra.Command{
	Use:   "list",
	Short: "List all valid tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, err := rootCmd.Flags().GetString("server-url")
		if err != nil {
			return err
		}
		client := api.NewClient(token, baseURL)
		var response struct {
			Data []models.Token `json:"data"`
		}
		data, err := client.Get("/tokens", nil)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &response); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
		tokens := response.Data

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tHashedToken\tCreatedAt\tExpiresAt\t")
		for _, t := range tokens {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t\n", t.ID, t.HashedToken, t.CreatedAt.Format(time.DateTime), t.ExpiresAt.Format(time.DateTime))
		}

		w.Flush()
		return nil
	},
}

func init() {
	tokensCmd.AddCommand(listTokensCmd)
	rootCmd.AddCommand(tokensCmd)
}
