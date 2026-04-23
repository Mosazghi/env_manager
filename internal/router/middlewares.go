package router

import (
	"env-manager/internal/handler"
	"env-manager/internal/repository"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func AuthRequired(tokenRepo *repository.TokenRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handler.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				handler.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header format must be Bearer <token>"})
				return
			}

			tokenString := parts[1]

			if len(tokenString) < 8 {
				handler.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
				return
			}

			validTokens, err := (*tokenRepo).FindAllValid(tokenString[:8])

			if err != nil {
				handler.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
				return
			}
			for _, token := range validTokens {
				if bcrypt.CompareHashAndPassword([]byte(token.HashedToken), []byte(tokenString)) == nil {
					next.ServeHTTP(w, r)
					return
				}
			}

			handler.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		})
	}
}
