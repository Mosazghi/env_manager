package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientGetBuildsAPIRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/projects" {
			t.Fatalf("expected path /api/projects, got %s", r.URL.Path)
		}

		if got := r.Header.Get("Authorization"); got != "Bearer token-123" {
			t.Fatalf("unexpected authorization header: %s", got)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := NewClient("token-123", server.URL)
	body, err := client.Get("/projects")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	if string(body) != `{"ok":true}` {
		t.Fatalf("unexpected response body: %s", string(body))
	}
}

func TestClientPostSendsJSONBody(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed reading body: %v", err)
		}

		var got payload
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("invalid json body: %v", err)
		}

		if got.Name != "demo" {
			t.Fatalf("expected payload name demo, got %s", got.Name)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"created"}`))
	}))
	defer server.Close()

	client := NewClient("token-123", server.URL)
	_, err := client.Post("/projects", payload{Name: "demo"})
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
}

func TestClientReturnsErrorForHTTPFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"bad request"}`))
	}))
	defer server.Close()

	client := NewClient("token-123", server.URL)
	_, err := client.Get("/projects")
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}

	if !strings.Contains(err.Error(), "bad request") {
		t.Fatalf("expected error to include response body, got: %v", err)
	}
}
