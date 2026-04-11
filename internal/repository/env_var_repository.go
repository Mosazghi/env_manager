package repository

import (
	"env-manager/internal/models"

	"gorm.io/gorm"
)

type EnvVarRepository interface {
	FindAll(page, limit int) ([]models.EnvVar, int64, error)
	FindByProjectID(projectID uint) ([]*models.EnvVar, error)
	FindByID(id uint) (*models.EnvVar, error)
	Create(envVar *models.EnvVar) error
	Update(envVar *models.EnvVar) error
	Delete(id uint) error
}

type envVarRepository struct {
	db *gorm.DB
}

func NewEnvVarRepository(db *gorm.DB) EnvVarRepository {
	return &envVarRepository{db}
}

func (r *envVarRepository) FindAll(page, limit int) ([]models.EnvVar, int64, error) {
	var envVars []models.EnvVar
	var total int64

	offset := (page - 1) * limit
	r.db.Model(&models.EnvVar{}).Count(&total)
	result := r.db.Preload("Project").Offset(offset).Limit(limit).Find(&envVars)

	return envVars, total, result.Error
}

func (r *envVarRepository) FindByProjectID(projectID uint) ([]*models.EnvVar, error) {
	var envVars []*models.EnvVar
	result := r.db.Where("project_id = ?", projectID).Find(&envVars)
	return envVars, result.Error
}

func (r *envVarRepository) FindByID(id uint) (*models.EnvVar, error) {
	var envVar models.EnvVar
	result := r.db.First(&envVar, id)
	return &envVar, result.Error
}

func (r *envVarRepository) Create(envVar *models.EnvVar) error {
	return r.db.Omit("Project").Create(envVar).Error
}

func (r *envVarRepository) Update(envVar *models.EnvVar) error {
	return r.db.Save(envVar).Error
}

func (r *envVarRepository) Delete(id uint) error {
	return r.db.Delete(&models.EnvVar{}, id).Error
}
