package router

import (
	"env-manager/internal/handler"
	"env-manager/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(projectHandler *handler.ProjectHandler, envVarHandler *handler.EnvVarHandlerand, tokenRepo *repository.TokenRepository) *gin.Engine {

	r := gin.Default()
	r.Use(AuthRequired(tokenRepo))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Users routes
	v1 := r.Group("/api")
	{
		// Users routes
		projects := v1.Group("/projects")
		{
			projects.GET("", projectHandler.GetAll)
			projects.GET("/:id", projectHandler.GetByID)
			projects.GET("/:id/env-vars", projectHandler.GetEnvVars)
			projects.POST("", projectHandler.Create)
			projects.PUT("/:id", projectHandler.Update)
			projects.DELETE("/:id", projectHandler.Delete)
		}

		// Env vars routes
		envVars := v1.Group("/env-vars")
		{
			envVars.GET("", envVarHandler.GetAll)
			envVars.GET("/:id", envVarHandler.FindByID)
			envVars.POST("", envVarHandler.Create)
			envVars.PUT("/:id", envVarHandler.Update)
			envVars.DELETE("/:id", envVarHandler.Delete)
		}
	}

	return r
}
