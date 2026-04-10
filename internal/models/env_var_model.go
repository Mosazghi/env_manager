package models

import "time"

type EnvVar struct {
	Project      Project `json:"-"    gorm:"foreignKey:ProjectID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Key          string `json:"key"`
	EncryptedVal string `json:"-" gorm:"column:encrypted_val"`
	Value        string `json:"value" gorm:"-"`
	ID           uint   `gorm:"primarykey"`
	ProjectID    int    `json:"project_id" gorm:"not null;index"`
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
