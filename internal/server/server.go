package server

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func NewServer(db *sql.DB) *gin.Engine {
	token := "4ebe17469d06d6823d9e9339ae97085d2c8bbca82f5e559ac3a48b6ecd7e8e67c20c2f35ef62c313e7eb752f42ff9525"

	router := gin.Default()
	router.Use(AuthRequired(token))

	router.POST("/env/projects/:project", CreateProject(db))

	router.GET("/env/projects", GetProjects(db))

	return router
}
