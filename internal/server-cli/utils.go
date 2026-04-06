package servercli

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseDuration(s string) time.Duration {
	if strings.HasSuffix(s, "d") {
		daysStr := strings.TrimSuffix(s, "d")
		if days, err := strconv.Atoi(daysStr); err == nil {
			return time.Duration(days) * 24 * time.Hour
		}
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		fmt.Printf("Invalid duration format: %s\n", s)
		os.Exit(1)
	}
	return d
}

func generateRandomToken() string {
	return rand.Text()
}

func getMasterPassphrase() (string, error) {
	credPath := os.Getenv("CREDENTIALS_DIRECTORY") + "/envm-passphrase"
	data, err := os.ReadFile(credPath)
	return strings.TrimSpace(string(data)), err
}
