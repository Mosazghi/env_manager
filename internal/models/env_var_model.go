package models

import "time"

type EnvVar struct {
	Project      Project   `json:"-" db:"-"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Key          string    `json:"key" db:"key"`
	EncryptedVal string    `json:"-" db:"encrypted_val"`
	Value        string    `json:"value" db:"-"`
	ID           uint      `db:"id"`
	ProjectID    int       `json:"project_id" db:"project_id"`
}

type CreateEnvVarRequest struct {
	Key       string `json:"key" binding:"required"`
	Value     string `json:"value" binding:"required"`
	ProjectID int    `json:"project_id" binding:"required"`
}

type UpdateEnvVarRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
