package repository

import (
	"env-manager/internal/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type TokenRepository interface {
	Create(token *models.Token) error
	FindAllValid(prefix string) ([]models.Token, error)
	DeleteExpired() error
}

type tokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) TokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) Create(token *models.Token) error {
	createdAt := token.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	} else {
		createdAt = createdAt.UTC()
	}

	expiresAt := token.ExpiresAt
	if expiresAt.IsZero() {
		return fmt.Errorf("token expiry is required")
	}
	expiresAt = expiresAt.UTC()

	result, err := r.db.Exec(
		`INSERT INTO tokens (prefix, hashed_token, created_at, expires_at) VALUES (?, ?, ?, ?)`,
		token.Prefix,
		token.HashedToken,
		createdAt,
		expiresAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	token.ID = uint(id)
	token.CreatedAt = createdAt.Local()
	token.ExpiresAt = expiresAt.Local()
	return nil
}

func (r *tokenRepository) DeleteExpired() error {
	_, err := r.db.Exec(`DELETE FROM tokens WHERE expires_at <= ?`, time.Now().UTC())
	return err
}

func (r *tokenRepository) FindAllValid(prefix string) ([]models.Token, error) {
	var tokens []models.Token
	err := r.db.Select(&tokens,
		`SELECT id, prefix, hashed_token, created_at, expires_at FROM tokens WHERE prefix = ? AND expires_at > ? ORDER BY id`,
		prefix,
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}
	return tokens, err
}
