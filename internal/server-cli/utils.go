package servercli

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseDuration(s string) (time.Duration, error) {
	if daysStr, ok := strings.CutSuffix(s, "d"); ok {
		if days, err := strconv.Atoi(daysStr); err == nil {
			return time.Duration(days) * 24 * time.Hour, nil
		}
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %s", s)
	}

	return d, nil
}

func generateRandomToken() string {
	return rand.Text()
}

func getMasterPassphrase() (string, error) {
	credPath := os.Getenv("CREDENTIALS_DIRECTORY") + "/envm-passphrase"
	data, err := os.ReadFile(credPath)

	return strings.TrimSpace(string(data)), err
}
