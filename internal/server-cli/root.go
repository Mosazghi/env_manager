package servercli

import (
	"env-manager/internal/config"
	"env-manager/internal/database"
	"env-manager/internal/models"
	"env-manager/internal/repository"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var rootCmd = &cobra.Command{
	Use:   "envm-server",
	Short: "env-manager Server CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage API tokens",
}

var tokenCreateCmd = &cobra.Command{
	Use:   "create [expires-in]",
	Short: "Create a new API token",
	Args:  cobra.NoArgs,
	RunE: func(servercli *cobra.Command, args []string) error {
		expiresIn, err := servercli.Flags().GetString("expires-in")
		if err != nil {
			return fmt.Errorf("invalid expires-in value: %w", err)
		}

		cfg := config.Load()

		db, err := database.NewSQLite(cfg.DBPath)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		c, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get raw DB connection: %w", err)
		}

		defer c.Close()
		tokenRepo := repository.NewTokenRepository(db)

		var token models.Token
		token.ExpiresAt = time.Now().Add(parseDuration(expiresIn))

		rawToken := generateRandomToken()
		hashedToken, err := bcrypt.GenerateFromPassword([]byte(rawToken), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash token: %w", err)
		}

		token.HashedToken = string(hashedToken)
		token.Prefix = rawToken[:8]

		if err := tokenRepo.Create(&token); err != nil {
			return fmt.Errorf("failed to create token: %w", err)
		}

		// silently delete expired tokens on each creation to avoid cluttering the database with old tokens
		tokenRepo.DeleteExpired()

		fmt.Printf("Token created: %s (expires in %s)\nCopy this token now, it won't be shown again!\n", rawToken, expiresIn)

		return nil
	},
}

var serverExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "Start the env-manager server",
	RunE: func(servercli *cobra.Command, args []string) error {
		svcConfig := &service.Config{
			Name:         "EnvManagerServer",
			DisplayName:  "Env Manager Server API Service",
			Description:  "Background HTTP API server for Env manager",
			Dependencies: serviceDependencies,
			Option: service.KeyValue{
				"OnFailure":      "restart",
				"StateDirectory": "envm",
			},
			Arguments: []string{"exec"},
		}

		prg := &program{}
		s, err := service.New(prg, svcConfig)
		if err != nil {
			log.Fatal(err)
			return err
		}

		if len(args) > 0 {
			action := args[0]
			if err := service.Control(s, action); err != nil {
				log.Fatalf("Valid actions: %q\nError: %v", service.ControlAction, err)
			}
			return nil
		}

		if err := s.Run(); err != nil {
			log.Fatal(err)
			return err

		}
		return nil
	}}

func init() {
	tokenCreateCmd.Flags().StringP("expires-in", "e", "1h", "Duration until the token expires (e.g. 30m, 2h, 10d)")
	tokenCmd.AddCommand(tokenCreateCmd)
	rootCmd.AddCommand(tokenCmd, serverExecCmd)
}
