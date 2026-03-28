package models

import "gorm.io/gorm"

type EnvVar struct {
	gorm.Model
	ProjectID int     `json:"project_id" gorm:"not null;index"`
	Project   Project `json:"-"    gorm:"foreignKey:ProjectID"`
	Key       string  `json:"key"`
	Value     string  `json:"value"`
}

type CreateEnvVarRequest struct {
	ProjectID int    `json:"project_id" binding:"required"`
	Key       string `json:"key" binding:"required"`
	Value     string `json:"value" binding:"required"`
}

type UpdateEnvVarRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}
