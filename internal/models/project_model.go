package models

import "time"

type Project struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description" gorm:"type:text"`
	ID          uint   `gorm:"primarykey"`
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
