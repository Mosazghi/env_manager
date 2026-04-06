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
	// systemd complicance for linux

	if runtime.GOOS == "linux" {
		return filepath.Join("/var/lib", "envm", "envm.db")
	}

	if stateDir := os.Getenv("STATE_DIRECTORY"); stateDir != "" {
		return filepath.Join(stateDir, "envm.db")
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("cannot determine config directory: ", err)
	}

	return filepath.Join(configDir, "envm", "envm.db")
}
