package repository

import (
	"database/sql"
	"env-manager/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
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
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) ProjectRepository {
	return &projectRepository{db}
}

func (r *projectRepository) FindAll() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Select(&projects, `SELECT id, name, description, created_at, updated_at FROM projects ORDER BY id`)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.Get(&project, `SELECT id, name, description, created_at, updated_at FROM projects WHERE id = ? LIMIT 1`, id)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) FindEnvVarsByID(id uint) ([]models.EnvVar, error) {
	var envVars []models.EnvVar
	err := r.db.Select(&envVars, `SELECT id, project_id, key, encrypted_val, created_at, updated_at FROM env_vars WHERE project_id = ? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	return envVars, nil
}

func (r *projectRepository) Create(project *models.Project) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`INSERT INTO projects (name, description, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		project.Name,
		project.Description,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	project.ID = uint(id)
	project.CreatedAt = now.Local()
	project.UpdatedAt = now.Local()
	return nil
}

func (r *projectRepository) Update(project *models.Project) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`UPDATE projects SET name = ?, description = ?, updated_at = ? WHERE id = ?`,
		project.Name,
		project.Description,
		time.Now(),
		project.ID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	project.UpdatedAt = now.Local()
	return nil
}

func (r *projectRepository) Delete(id uint) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}
