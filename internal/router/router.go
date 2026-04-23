package router

import (
	"env-manager/internal/handler"
	"env-manager/internal/repository"
	"net/http"
)

func Setup(projectHandler *handler.ProjectHandler, envVarHandler *handler.EnvVarHandlerand, tokenRepo *repository.TokenRepository) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		handler.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/projects", projectHandler.GetAll)
	mux.HandleFunc("GET /api/projects/", projectHandler.GetAll)
	mux.HandleFunc("GET /api/projects/{id}", projectHandler.GetByID)
	mux.HandleFunc("GET /api/projects/{id}/env-vars", projectHandler.GetEnvVars)
	mux.HandleFunc("POST /api/projects", projectHandler.Create)
	mux.HandleFunc("POST /api/projects/", projectHandler.Create)
	mux.HandleFunc("PUT /api/projects/{id}", projectHandler.Update)
	mux.HandleFunc("DELETE /api/projects/{id}", projectHandler.Delete)

	mux.HandleFunc("GET /api/env-vars", envVarHandler.GetAll)
	mux.HandleFunc("GET /api/env-vars/", envVarHandler.GetAll)
	mux.HandleFunc("GET /api/env-vars/{id}", envVarHandler.FindByID)
	mux.HandleFunc("POST /api/env-vars", envVarHandler.Create)
	mux.HandleFunc("POST /api/env-vars/", envVarHandler.Create)
	mux.HandleFunc("PUT /api/env-vars/{id}", envVarHandler.Update)
	mux.HandleFunc("DELETE /api/env-vars/{id}", envVarHandler.Delete)

	return AuthRequired(tokenRepo)(mux)
}
