package router

import (
	"env-manager/internal/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(projectHandler *handler.ProjectHandler, envVarHandler *handler.EnvVarHandlerand) *gin.Engine {
	token := "4ebe17469d06d6823d9e9339ae97085d2c8bbca82f5e559ac3a48b6ecd7e8e67c20c2f35ef62c313e7eb752f42ff9525"

	r := gin.Default()
	r.Use(AuthRequired(token))

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
