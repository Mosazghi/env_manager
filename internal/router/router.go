package router

import (
	"env-manager/internal/handler"
	"env-manager/internal/repository"
	"net/http"
)

type Router struct {
	ProjectHandler *handler.ProjectHandler
	EnvVarHandler  *handler.EnvVarHandlerand
	TokenHandler   *handler.TokenHandler
	TokenRepo      *repository.TokenRepository
}

func Setup(router *Router) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		handler.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/projects", router.ProjectHandler.GetAll)
	mux.HandleFunc("GET /api/projects/", router.ProjectHandler.GetAll)
	mux.HandleFunc("GET /api/projects/{id}", router.ProjectHandler.GetByID)
	mux.HandleFunc("GET /api/projects/{id}/env-vars", router.ProjectHandler.GetEnvVars)
	mux.HandleFunc("POST /api/projects", router.ProjectHandler.Create)
	mux.HandleFunc("POST /api/projects/", router.ProjectHandler.Create)
	mux.HandleFunc("PUT /api/projects/{id}", router.ProjectHandler.Update)
	mux.HandleFunc("DELETE /api/projects/{id}", router.ProjectHandler.Delete)

	mux.HandleFunc("GET /api/env-vars", router.EnvVarHandler.GetAll)
	mux.HandleFunc("GET /api/env-vars/", router.EnvVarHandler.GetAll)
	mux.HandleFunc("GET /api/env-vars/{id}", router.EnvVarHandler.FindByID)
	mux.HandleFunc("POST /api/env-vars", router.EnvVarHandler.Create)
	mux.HandleFunc("POST /api/env-vars/", router.EnvVarHandler.Create)
	mux.HandleFunc("PUT /api/env-vars/{id}", router.EnvVarHandler.Update)
	mux.HandleFunc("DELETE /api/env-vars/{id}", router.EnvVarHandler.Delete)

	mux.HandleFunc("GET /api/tokens", router.TokenHandler.GetAll)

	return AuthRequired(router.TokenRepo)(mux)
}
