package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
)

type Project struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func getMasterPassphrase() (string, error) {
	credPath := os.Getenv("CREDENTIALS_DIRECTORY") + "/envm-passphrase"
	data, err := os.ReadFile(credPath)
	return strings.TrimSpace(string(data)), err
}

func GenerateAuthToken(passphrase string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(passphrase))
	h.Write(salt)

	token := append(salt, h.Sum(nil)...)

	return hex.EncodeToString(token), nil
}

func AuthRequired(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			return
		}

		tokenString := parts[1]

		if tokenString != expectedToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}

func DBInit() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./envs.db")
	if err != nil {
		return nil, err
	}

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS project (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT,
	    	created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	log.Println("Table 'project' created successfully")
	return db, nil
}

func main() {
	token := "4ebe17469d06d6823d9e9339ae97085d2c8bbca82f5e559ac3a48b6ecd7e8e67c20c2f35ef62c313e7eb752f42ff9525"

	fmt.Println("TOKEN: ", token)
	db, err := DBInit()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(AuthRequired(token))

	router.POST("/env/projects/:project", func(c *gin.Context) {
		projectName := c.Param("project")
		_, err := db.Exec(`INSERT INTO project (name) VALUES (?)`, projectName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Project created", "project": projectName})
	})

	router.GET("/env/projects", func(c *gin.Context) {
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
	})

	router.Run(":8080")
}
