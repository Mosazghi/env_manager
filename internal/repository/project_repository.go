package repository

import (
	"env-manager/internal/models"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	FindAll(page, limit int) ([]models.Project, int64, error)
	FindByID(id uint) (*models.Project, error)
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

func (r *projectRepository) FindAll(page, limit int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	offset := (page - 1) * limit
	r.db.Model(&models.Project{}).Count(&total)
	result := r.db.Offset(offset).Limit(limit).Find(&projects)
	return projects, total, result.Error
}

func (r *projectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	result := r.db.First(&project, id)
	return &project, result.Error
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
