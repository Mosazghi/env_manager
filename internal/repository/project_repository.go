package repository

import (
	"env-manager/internal/models"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	FindAll() ([]models.Project, error)
	FindByID(id uint) (*models.Project, error)
	FindEnvVarsByID(id uint) ([]models.EnvVar, error)
	Create(project *models.Project) error
	Update(project *models.Project) error
	Delete(id uint) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db}
}

func (r *projectRepository) FindAll() ([]models.Project, error) {
	var projects []models.Project

	result := r.db.Find(&projects)
	return projects, result.Error
}

func (r *projectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	result := r.db.First(&project, id)
	return &project, result.Error
}

func (r *projectRepository) FindEnvVarsByID(id uint) ([]models.EnvVar, error) {
	var envVars []models.EnvVar
	result := r.db.Where("project_id = ?", id).Find(&envVars)
	return envVars, result.Error
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}
