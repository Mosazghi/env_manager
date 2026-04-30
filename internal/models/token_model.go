package models

import (
	"time"
)

type Token struct {
	ExpiresAt   time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	HashedToken string    `json:"hashedToken" db:"hashed_token"`
	// Only the 8 first chars of the token for faster lookup
	Prefix string `json:"-" db:"prefix"`
	ID     uint   `json:"id" db:"id"`
}
