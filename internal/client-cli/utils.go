package clientcli

import (
	"bufio"
	"os"
	"strings"
)

func truncateWithEllipsis(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	if maxLength <= 3 {
		return s[:maxLength]
	}
	return s[:maxLength-3] + "..."
}

func truncateProjectDescription(desc string) string {
	return truncateWithEllipsis(desc, 30)
}

func truncateProjectName(name string) string {
	return truncateWithEllipsis(name, 20)
}

func getLocalEnvVars(filePath string) (map[string]string, bool) {
	envVars := make(map[string]string)

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o644)
	if err != nil {
		return envVars, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		variable := strings.SplitN(line, "=", 2)
		if len(variable) != 2 {
			continue
		}
		key := variable[0]
		val := variable[1]

		isValid := key != "" && !strings.HasPrefix(key, "#") && val != ""
		if isValid {
			envVars[key] = val
		}
	}

	return envVars, true
}
