package clientcli

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectCommandsRunE(t *testing.T) {
	listCalled := false
	createCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/projects":
			listCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":[{"id":1,"name":"demo","description":"test","createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/projects/":
			createCalled = true
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"data":{}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	rootCmd.SetArgs([]string{"--token", "token-123", "--server-url", server.URL, "projects", "list"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("projects list returned error: %v", err)
	}

	rootCmd.SetArgs([]string{"--token", "token-123", "--server-url", server.URL, "projects", "create", "demo", "test"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("projects create returned error: %v", err)
	}

	if !listCalled || !createCalled {
		t.Fatalf("expected both list and create endpoints to be called, list=%v create=%v", listCalled, createCalled)
	}
}

func TestEnvVarCreateAndLoadCommandsRunE(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change cwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWd)
	})

	createCalled := false
	loadCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/env-vars":
			createCalled = true
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"data":{}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/projects/1/env-vars":
			loadCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":[{"key":"API_KEY","value":"secret"}]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	rootCmd.SetArgs([]string{"--token", "token-123", "--server-url", server.URL, "--project-id", "1", "env-vars", "create", "--key", "API_KEY", "--value", "secret"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("env-vars create returned error: %v", err)
	}

	rootCmd.SetArgs([]string{"--token", "token-123", "--server-url", server.URL, "--project-id", "1", "env-vars", "load"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("env-vars load returned error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, ".env"))
	if err != nil {
		t.Fatalf("failed reading generated .env file: %v", err)
	}

	if !strings.Contains(string(content), "API_KEY=secret") {
		t.Fatalf("expected generated .env to contain API_KEY, got %q", string(content))
	}

	if !createCalled || !loadCalled {
		t.Fatalf("expected create and load endpoints to be called, create=%v load=%v", createCalled, loadCalled)
	}
}

func TestEnvVarSyncCommandRunEForceUpdate(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(filePath, []byte("NEW_KEY=new\nSHARED_KEY=local\n"), 0o644); err != nil {
		t.Fatalf("failed writing local env file: %v", err)
	}

	getCalled := false
	postCalled := false
	putCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/projects/1/env-vars":
			getCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":[{"id":5,"key":"SHARED_KEY","value":"remote"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/env-vars":
			postCalled = true
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"data":{}}`))
		case r.Method == http.MethodPut && r.URL.Path == "/api/env-vars/5":
			putCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	rootCmd.SetArgs([]string{"--token", "token-123", "--server-url", server.URL, "--project-id", "1", "env-vars", "sync", "--force-update", "--file-path", filePath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("env-vars sync returned error: %v", err)
	}

	if !getCalled || !postCalled || !putCalled {
		t.Fatalf("expected GET/POST/PUT calls, got get=%v post=%v put=%v", getCalled, postCalled, putCalled)
	}
}

func TestStoreTokenAndServerURL(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "envm.global.json")

	if err := storeToken("token-123", filePath); err != nil {
		t.Fatalf("storeToken returned error: %v", err)
	}

	if err := storeServerURL("http://localhost:8080", filePath); err != nil {
		t.Fatalf("storeServerURL returned error: %v", err)
	}

	config, err := getGlobalConfigData(filePath)
	if err != nil {
		t.Fatalf("getGlobalConfigData returned error: %v", err)
	}

	if config.Token != "token-123" {
		t.Fatalf("expected token token-123, got %q", config.Token)
	}

	if config.ServerURL != "http://localhost:8080" {
		t.Fatalf("expected serverURL http://localhost:8080, got %q", config.ServerURL)
	}
}

func TestSetProjectID(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), ".envm.local.json")

	if err := setProjectID("42", filePath); err != nil {
		t.Fatalf("setProjectID returned error: %v", err)
	}

	config, err := getLocalConfigData(filePath)
	if err != nil {
		t.Fatalf("getLocalConfigData returned error: %v", err)
	}

	if config.ProjectID != "42" {
		t.Fatalf("expected projectID 42, got %q", config.ProjectID)
	}
}

func TestGetStoredValuesFromDefaultPaths(t *testing.T) {
	configHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)

	globalConfigDir := filepath.Join(configHome, "envm")
	if err := os.MkdirAll(globalConfigDir, 0o755); err != nil {
		t.Fatalf("failed to create global config directory: %v", err)
	}

	globalConfigPath := filepath.Join(globalConfigDir, "envm.global.json")
	globalConfigContent := `{"token":"token-from-file","serverUrl":"http://example.com"}`
	if err := os.WriteFile(globalConfigPath, []byte(globalConfigContent), 0o644); err != nil {
		t.Fatalf("failed to write global config file: %v", err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	workingDir := t.TempDir()
	if err := os.Chdir(workingDir); err != nil {
		t.Fatalf("failed to change cwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWd)
	})

	localConfigPath := filepath.Join(workingDir, ".envm.local.json")
	if err := os.WriteFile(localConfigPath, []byte(`{"projectID":"99"}`), 0o644); err != nil {
		t.Fatalf("failed to write local config file: %v", err)
	}

	token, err := GetStoredToken()
	if err != nil {
		t.Fatalf("GetStoredToken returned error: %v", err)
	}

	url, err := GetStoredServerURL()
	if err != nil {
		t.Fatalf("GetStoredServerURL returned error: %v", err)
	}

	pID, err := GetStoredProjectID()
	if err != nil {
		t.Fatalf("GetStoredProjectID returned error: %v", err)
	}

	if token != "token-from-file" {
		t.Fatalf("expected token token-from-file, got %q", token)
	}

	if url != "http://example.com" {
		t.Fatalf("expected serverURL http://example.com, got %q", url)
	}

	if pID != "99" {
		t.Fatalf("expected projectID 99, got %q", pID)
	}
}
