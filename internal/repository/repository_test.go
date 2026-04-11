package repository

import (
	"path/filepath"
	"testing"
	"time"

	"env-manager/internal/database"
	"env-manager/internal/models"
)

func newTestDB(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "repo.db")
}

func TestProjectRepositoryCRUD(t *testing.T) {
	db, err := database.NewSQLite(newTestDB(t))
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	repo := NewProjectRepository(db)

	project := &models.Project{Name: "alpha", Description: "first project"}
	if err := repo.Create(project); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	found, err := repo.FindByID(project.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}

	if found.Name != "alpha" {
		t.Fatalf("expected project name alpha, got %s", found.Name)
	}

	projects, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll returned error: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	found.Description = "updated"
	if err := repo.Update(found); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	updated, err := repo.FindByID(project.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}

	if updated.Description != "updated" {
		t.Fatalf("expected updated description, got %s", updated.Description)
	}

	if err := repo.Delete(project.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	if _, err := repo.FindByID(project.ID); err == nil {
		t.Fatal("expected error after deleting project")
	}
}

func TestProjectRepositoryFindEnvVarsByID(t *testing.T) {
	db, err := database.NewSQLite(newTestDB(t))
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	projectRepo := NewProjectRepository(db)
	envRepo := NewEnvVarRepository(db)

	project := &models.Project{Name: "beta", Description: "with env vars"}
	if err := projectRepo.Create(project); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	envVar := &models.EnvVar{ProjectID: int(project.ID), Key: "PORT", EncryptedVal: "enc"}
	if err := envRepo.Create(envVar); err != nil {
		t.Fatalf("Create env var returned error: %v", err)
	}

	envs, err := projectRepo.FindEnvVarsByID(project.ID)
	if err != nil {
		t.Fatalf("FindEnvVarsByID returned error: %v", err)
	}

	if len(envs) != 1 {
		t.Fatalf("expected 1 env var, got %d", len(envs))
	}

	if envs[0].Key != "PORT" {
		t.Fatalf("expected env key PORT, got %s", envs[0].Key)
	}
}

func TestEnvVarRepositoryCRUDAndPagination(t *testing.T) {
	db, err := database.NewSQLite(newTestDB(t))
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	projectRepo := NewProjectRepository(db)
	repo := NewEnvVarRepository(db)

	project := &models.Project{Name: "gamma", Description: "env vars"}
	if err := projectRepo.Create(project); err != nil {
		t.Fatalf("Create project returned error: %v", err)
	}

	envA := &models.EnvVar{ProjectID: int(project.ID), Key: "A", EncryptedVal: "enc-a"}
	envB := &models.EnvVar{ProjectID: int(project.ID), Key: "B", EncryptedVal: "enc-b"}
	if err := repo.Create(envA); err != nil {
		t.Fatalf("Create envA returned error: %v", err)
	}
	if err := repo.Create(envB); err != nil {
		t.Fatalf("Create envB returned error: %v", err)
	}

	page, total, err := repo.FindAll(1, 1)
	if err != nil {
		t.Fatalf("FindAll returned error: %v", err)
	}

	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}

	if len(page) != 1 {
		t.Fatalf("expected page size 1, got %d", len(page))
	}

	byProject, err := repo.FindByProjectID(project.ID)
	if err != nil {
		t.Fatalf("FindByProjectID returned error: %v", err)
	}

	if len(byProject) != 2 {
		t.Fatalf("expected 2 env vars for project, got %d", len(byProject))
	}

	found, err := repo.FindByID(envA.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}

	found.Key = "A_UPDATED"
	if err := repo.Update(found); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	updated, err := repo.FindByID(envA.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}

	if updated.Key != "A_UPDATED" {
		t.Fatalf("expected updated key, got %s", updated.Key)
	}

	if err := repo.Delete(envA.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	if _, err := repo.FindByID(envA.ID); err == nil {
		t.Fatal("expected error after deleting env var")
	}
}

func TestTokenRepositoryFindAllValidAndDeleteExpired(t *testing.T) {
	db, err := database.NewSQLite(newTestDB(t))
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	repo := NewTokenRepository(db)
	now := time.Now()

	validToken := &models.Token{
		Prefix:      "abcdefgh",
		HashedToken: "valid-hash",
		CreatedAt:   now,
		ExpiresAt:   now.Add(time.Hour),
	}
	expiredToken := &models.Token{
		Prefix:      "abcdefgh",
		HashedToken: "expired-hash",
		CreatedAt:   now,
		ExpiresAt:   now.Add(-time.Hour),
	}

	if err := repo.Create(validToken); err != nil {
		t.Fatalf("Create valid token returned error: %v", err)
	}

	if err := repo.Create(expiredToken); err != nil {
		t.Fatalf("Create expired token returned error: %v", err)
	}

	valid, err := repo.FindAllValid("abcdefgh")
	if err != nil {
		t.Fatalf("FindAllValid returned error: %v", err)
	}

	if len(valid) != 1 {
		t.Fatalf("expected 1 valid token, got %d", len(valid))
	}

	if valid[0].HashedToken != "valid-hash" {
		t.Fatalf("expected valid token hash, got %s", valid[0].HashedToken)
	}

	if err := repo.DeleteExpired(); err != nil {
		t.Fatalf("DeleteExpired returned error: %v", err)
	}

	var count int64
	if err := db.Model(&models.Token{}).Count(&count).Error; err != nil {
		t.Fatalf("failed counting tokens: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected only 1 token to remain, got %d", count)
	}
}
