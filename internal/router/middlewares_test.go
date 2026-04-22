package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"env-manager/internal/models"
	"env-manager/internal/repository"
)

type tokenRepoErrorMock struct{}

func (tokenRepoErrorMock) Create(token *models.Token) error { return nil }

func (tokenRepoErrorMock) FindAllValid(prefix string) ([]models.Token, error) {
	return nil, errors.New("query failed")
}

func (tokenRepoErrorMock) DeleteExpired() error { return nil }

func TestAuthRequiredHandlesRepositoryError(t *testing.T) {
	repoIface := repository.TokenRepository(tokenRepoErrorMock{})
	secured := AuthRequired(&repoIface)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/secured", nil)
	req.Header.Set("Authorization", "Bearer abcdefgh-token")
	w := httptest.NewRecorder()
	secured.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d with body %s", w.Code, w.Body.String())
	}
}
