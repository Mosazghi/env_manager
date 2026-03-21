package server

import (
	"database/sql"
	"net/http"

	. "env-manager/internal/shared"

	"github.com/gin-gonic/gin"
)

func CreateProject(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectName := c.Param("project")
		row := db.QueryRow("SELET name FROM project WHERE name = (?)", projectName)
		if row != nil {
			c.JSON(http.StatusOK, gin.H{"message": "Project already exists!", "projectName": projectName})
			return
		}
		_, err := db.Exec(`INSERT INTO project (name) VALUES (?)`, projectName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Project created", "project": projectName})
	}
}

func GetProjects(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Convert stored timestamps to local time for API responses.
		rows, err := db.Query("SELECT id, name, datetime(created_at, 'localtime') FROM project")
		var projects []Project
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var p Project
			if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			projects = append(projects, p)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, projects)
	}
}
