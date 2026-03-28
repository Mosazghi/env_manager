package main

import (
	"context"
	"env-manager/internal/config"
	"env-manager/internal/database"
	"env-manager/internal/handler"
	"env-manager/internal/repository"
	"env-manager/internal/router"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/kardianos/service"
)

func getMasterPassphrase() (string, error) {
	credPath := os.Getenv("CREDENTIALS_DIRECTORY") + "/envm-passphrase"
	data, err := os.ReadFile(credPath)
	return strings.TrimSpace(string(data)), err
}

var logger service.Logger

// program holds the HTTP server so Stop() can shut it down cleanly.
type program struct {
	srv *http.Server
}

func (p *program) Start(s service.Service) error {
	cfg := config.Load()
	// Start() must be non-blocking — spin the server in a goroutine.
	db, err := database.NewSQLite(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	projectRepo := repository.NewProjectRepository(db)
	envVarRepo := repository.NewEnvVarRepository(db)

	// Wire up handlers
	projectHandler := handler.NewProjectHandler(projectRepo)
	envVarHandler := handler.NewEnvVarHandler(projectRepo, envVarRepo)

	router := router.Setup(projectHandler, envVarHandler)
	p.srv = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
	go func() {
		if err := p.srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Errorf("HTTP server error: %v", err)
		}
	}()
	logger.Infof("Service started on port %s", cfg.Port)
	return nil
}

func (p *program) Stop(s service.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.srv.Shutdown(ctx); err != nil {
		return err
	}
	logger.Infof("Service stopped")
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:         "EnvManagerServer",
		DisplayName:  "Env Manager Server API Service",
		Description:  "Background HTTP API server for MyApp",
		Dependencies: []string{"After=network.target"},
		Option:       service.KeyValue{"OnFailure": "restart"},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		action := os.Args[1]
		if err := service.Control(s, action); err != nil {
			log.Fatalf("Valid actions: %q\nError: %v", service.ControlAction, err)
		}
		return
	}

	if err := s.Run(); err != nil {
		logger.Error(err)
	}
}
