package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	DBPath string
	Env    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	return &Config{
		Port:   getEnv("PORT", "8080"),
		DBPath: getEnv("DB_PATH", defaultDBPath()),
		Env:    getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func defaultDBPath() string {
	switch runtime.GOOS {
	case "linux":
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			return filepath.Join(xdg, "envm", "envm.db")
		}
		return filepath.Join(os.Getenv("HOME"), ".local", "share", "envm", "envm.db")

	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "envm", "envm.db")

	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "envm", "envm.db")

	default:
		return filepath.Join(os.Getenv("HOME"), ".envm", "envm.db")
	}
}
