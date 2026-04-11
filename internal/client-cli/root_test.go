package clientcli

import (
	"testing"
)

func TestGenerateStars(t *testing.T) {
	if got := generateStars("abcd"); got != "****" {
		t.Fatalf("expected ****, got %q", got)
	}

	if got := generateStars(""); got != "" {
		t.Fatalf("expected empty output for empty input, got %q", got)
	}
}
