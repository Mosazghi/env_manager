package repository

import (
	"database/sql"
	"env-manager/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
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
	db *sqlx.DB
}

func NewEnvVarRepository(db *sqlx.DB) EnvVarRepository {
	return &envVarRepository{db}
}

func (r *envVarRepository) FindAll(page, limit int) ([]models.EnvVar, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	var total int64
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM env_vars`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var envVars []models.EnvVar
	err := r.db.Select(&envVars, `SELECT id, project_id, key, encrypted_val, created_at, updated_at FROM env_vars ORDER BY id LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return envVars, total, nil
}

func (r *envVarRepository) FindByProjectID(projectID uint) ([]*models.EnvVar, error) {
	var envVars []*models.EnvVar
	err := r.db.Select(&envVars, `SELECT id, project_id, key, encrypted_val, created_at, updated_at FROM env_vars WHERE project_id = ? ORDER BY id`, projectID)
	if err != nil {
		return nil, err
	}
	return envVars, nil
}
func (r *envVarRepository) FindByID(id uint) (*models.EnvVar, error) {
	var envVar models.EnvVar
	err := r.db.Get(&envVar, `SELECT id, project_id, key, encrypted_val, created_at, updated_at FROM env_vars WHERE id = ? LIMIT 1`, id)
	if err != nil {
		return nil, err
	}
	return &envVar, nil

}

func (r *envVarRepository) Create(envVar *models.EnvVar) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`INSERT INTO env_vars (project_id, key, encrypted_val, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		envVar.ProjectID,
		envVar.Key,
		envVar.EncryptedVal,
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

	envVar.ID = uint(id)
	envVar.CreatedAt = now.Local()
	envVar.UpdatedAt = now.Local()
	return nil
}

func (r *envVarRepository) Update(envVar *models.EnvVar) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`UPDATE env_vars SET project_id = ?, key = ?, encrypted_val = ?, updated_at = ? WHERE id = ?`,
		envVar.ProjectID,
		envVar.Key,
		envVar.EncryptedVal,
		time.Now(),
		envVar.ID,
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

	envVar.UpdatedAt = now.Local()
	return nil
}

func (r *envVarRepository) Delete(id uint) error {
	_, err := r.db.Exec(`DELETE FROM env_vars WHERE id = ?`, id)
	return err
}
