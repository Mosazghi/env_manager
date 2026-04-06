package repository

import (
	"env-manager/internal/models"
	"time"

	"gorm.io/gorm"
)

type TokenRepository interface {
	Create(token *models.Token) error
	FindAllValid(prefix string) ([]models.Token, error)
	DeleteExpired() error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) Create(token *models.Token) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at <= ?", time.Now()).Delete(&models.Token{}).Error
}

func (r *tokenRepository) FindAllValid(prefix string) ([]models.Token, error) {
	var tokens []models.Token
	result := r.db.Where("prefix = ? AND expires_at > ?", prefix, time.Now()).Find(&tokens)
	return tokens, result.Error
}
