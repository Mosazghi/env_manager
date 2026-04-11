package api

import (
	"net/url"
)

func isValidPath(rawPath string) bool {
	u, err := url.Parse(rawPath)
	if err != nil {
		return false
	}

	return u.Host != "" && u.Path != "" && u.Scheme != ""
}
