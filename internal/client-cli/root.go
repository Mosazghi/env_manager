package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var token string

var rootCmd = &cobra.Command{
	Use:   "envm",
	Short: "env-manager CLI",
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
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("scanner encountered an error: %s", err)
	}

	rootCmd.PersistentFlags().StringVar(&token, "token", token, "API token (or set ENVM_TOKEN)")
}
