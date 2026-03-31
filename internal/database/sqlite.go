package database

import (
	"os"
	"path/filepath"

	"env-manager/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSQLite(path string) (*gorm.DB, error) {
	if err := ensureDir(path); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(&models.Project{}, &models.EnvVar{}); err != nil {
		return nil, err
	}

	return db, nil
}

func ensureDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0o700)
}
