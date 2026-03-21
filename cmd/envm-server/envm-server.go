package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"env-manager/internal/db"
	"env-manager/internal/server"

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
	// Start() must be non-blocking — spin the server in a goroutine.
	db, err := db.DBInit()
	if err != nil {
		panic(err)
	}

	router := server.NewServer(db)
	p.srv = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := p.srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Errorf("HTTP server error: %v", err)
		}
	}()
	logger.Info("Service started on :8080")
	return nil
}

func (p *program) Stop(s service.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.srv.Shutdown(ctx); err != nil {
		return err
	}
	logger.Info("Service stopped")
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:         "EnvManagerServer",
		DisplayName:  "My App HTTP Service",
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
		// Supported: install, uninstall, start, stop, restart
		if err := service.Control(s, action); err != nil {
			log.Fatalf("Valid actions: %q\nError: %v", service.ControlAction, err)
		}
		return
	}

	if err := s.Run(); err != nil {
		logger.Error(err)
	}
}
