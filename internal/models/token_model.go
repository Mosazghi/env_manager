package models

import (
	"time"
)

type Token struct {
	ExpiresAt   time.Time `json:"expiresAt" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"not null"`
	HashedToken string    `json:"hashedToken" gorm:"not null"`
	// Only the 8 first chars of the token for faster lookup
	Prefix string `json:"prefix"`
	ID     uint   `json:"id" gorm:"primarykey"`
}
