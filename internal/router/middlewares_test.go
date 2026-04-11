package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"env-manager/internal/models"
	"env-manager/internal/repository"

	"github.com/gin-gonic/gin"
)

type tokenRepoErrorMock struct{}

func (tokenRepoErrorMock) Create(token *models.Token) error { return nil }

func (tokenRepoErrorMock) FindAllValid(prefix string) ([]models.Token, error) {
	return nil, errors.New("query failed")
}

func (tokenRepoErrorMock) DeleteExpired() error { return nil }

func TestAuthRequiredHandlesRepositoryError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repoIface := repository.TokenRepository(tokenRepoErrorMock{})
	r := gin.New()
	r.Use(AuthRequired(&repoIface))
	r.GET("/secured", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/secured", nil)
	req.Header.Set("Authorization", "Bearer abcdefgh-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d with body %s", w.Code, w.Body.String())
	}
}
