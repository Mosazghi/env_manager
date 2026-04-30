package handler

import (
	"env-manager/internal/repository"
	"net/http"
)

type TokenHandler struct {
	repo repository.TokenRepository
}

func NewTokenHandler(repo repository.TokenRepository) *TokenHandler {
	return &TokenHandler{repo}
}

func (h *TokenHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tokens, err := h.repo.FindAllValid("*")
	for i := range tokens {
		len := len(tokens[i].HashedToken)
		if len > 8 {
			tokens[i].HashedToken = tokens[i].HashedToken[:4] + "..." + tokens[i].HashedToken[len-4:]
		}

	}
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	WriteJSON(w, http.StatusOK, ToResponse(true, "Tokens retrieved", tokens))
}
