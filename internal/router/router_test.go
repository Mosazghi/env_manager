package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"env-manager/internal/database"
	"env-manager/internal/handler"
	"env-manager/internal/models"
	"env-manager/internal/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type apiResponse[T any] struct {
	Sucess  bool   `json:"sucess"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func setupTestRouter(t *testing.T) (*gin.Engine, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	t.Setenv("ENVM_MASTER_KEY_FILE", filepath.Join(t.TempDir(), "master.key"))

	db, err := database.NewSQLite(filepath.Join(t.TempDir(), "router.db"))
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	projectRepo := repository.NewProjectRepository(db)
	envRepo := repository.NewEnvVarRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	rawToken := "abcdefgh-super-secret-token"
	hash, err := bcrypt.GenerateFromPassword([]byte(rawToken), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash token: %v", err)
	}

	if err := tokenRepo.Create(&models.Token{
		Prefix:      rawToken[:8],
		HashedToken: string(hash),
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(time.Hour),
	}); err != nil {
		t.Fatalf("failed to seed token: %v", err)
	}

	projectHandler := handler.NewProjectHandler(projectRepo)
	envHandler := handler.NewEnvVarHandler(projectRepo, envRepo)

	return Setup(projectHandler, envHandler, &tokenRepo), rawToken
}

func doJSONRequest(t *testing.T, r http.Handler, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthRequiredRejectsMissingHeader(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := doJSONRequest(t, r, http.MethodGet, "/health", "", nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d with body %s", w.Code, w.Body.String())
	}
}

func TestProjectAndEnvVarFlow(t *testing.T) {
	r, token := setupTestRouter(t)

	createProject := doJSONRequest(t, r, http.MethodPost, "/api/projects", token, map[string]string{
		"name":        "demo",
		"description": "demo project",
	})
	if createProject.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating project, got %d with body %s", createProject.Code, createProject.Body.String())
	}

	var createProjectResp apiResponse[models.Project]
	if err := json.Unmarshal(createProject.Body.Bytes(), &createProjectResp); err != nil {
		t.Fatalf("failed to unmarshal project response: %v", err)
	}

	projectID := createProjectResp.Data.ID
	if projectID == 0 {
		t.Fatal("expected non-zero project id")
	}

	createEnv := doJSONRequest(t, r, http.MethodPost, "/api/env-vars", token, map[string]any{
		"key":        "API_KEY",
		"value":      "top-secret",
		"project_id": int(projectID),
	})
	if createEnv.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating env var, got %d with body %s", createEnv.Code, createEnv.Body.String())
	}

	var createEnvResp apiResponse[models.EnvVar]
	if err := json.Unmarshal(createEnv.Body.Bytes(), &createEnvResp); err != nil {
		t.Fatalf("failed to unmarshal env-var response: %v", err)
	}

	envID := createEnvResp.Data.ID
	if envID == 0 {
		t.Fatal("expected non-zero env var id")
	}

	getProjectEnv := doJSONRequest(t, r, http.MethodGet, fmt.Sprintf("/api/projects/%d/env-vars", projectID), token, nil)
	if getProjectEnv.Code != http.StatusOK {
		t.Fatalf("expected 200 getting project env vars, got %d with body %s", getProjectEnv.Code, getProjectEnv.Body.String())
	}

	var projectEnvResp apiResponse[[]models.EnvVar]
	if err := json.Unmarshal(getProjectEnv.Body.Bytes(), &projectEnvResp); err != nil {
		t.Fatalf("failed to unmarshal project env response: %v", err)
	}

	if len(projectEnvResp.Data) != 1 {
		t.Fatalf("expected one env var, got %d", len(projectEnvResp.Data))
	}

	if projectEnvResp.Data[0].Value != "top-secret" {
		t.Fatalf("expected decrypted value top-secret, got %q", projectEnvResp.Data[0].Value)
	}

	getByID := doJSONRequest(t, r, http.MethodGet, fmt.Sprintf("/api/env-vars/%d", envID), token, nil)
	if getByID.Code != http.StatusOK {
		t.Fatalf("expected 200 getting env var by id, got %d with body %s", getByID.Code, getByID.Body.String())
	}

	var envByIDResp apiResponse[models.EnvVar]
	if err := json.Unmarshal(getByID.Body.Bytes(), &envByIDResp); err != nil {
		t.Fatalf("failed to unmarshal env-by-id response: %v", err)
	}

	if envByIDResp.Data.Value != "top-secret" {
		t.Fatalf("expected decrypted value top-secret, got %q", envByIDResp.Data.Value)
	}
}
