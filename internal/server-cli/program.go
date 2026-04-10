package servercli

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"env-manager/internal/config"
	"env-manager/internal/database"
	"env-manager/internal/handler"
	"env-manager/internal/repository"
	"env-manager/internal/router"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/kardianos/service"
)

type program struct {
	srv *http.Server
}

func (p *program) Start(s service.Service) error {
	cfg := config.Load()
	fmt.Printf("Loaded config: Port=%s, DBPath=%s, Env=%s\n", cfg.Port, cfg.DBPath, cfg.Env)
	db, err := database.NewSQLite(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	projectRepo := repository.NewProjectRepository(db)
	envVarRepo := repository.NewEnvVarRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	projectHandler := handler.NewProjectHandler(projectRepo)
	envVarHandler := handler.NewEnvVarHandler(projectRepo, envVarRepo)

	router := router.Setup(projectHandler, envVarHandler, &tokenRepo)
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
