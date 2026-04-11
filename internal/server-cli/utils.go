package servercli

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
