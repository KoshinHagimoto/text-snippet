package config

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGenerateSignedURL(t *testing.T) {
	os.Setenv("SIGNED_URL_SECRET", "test_secret_key")

	userID := 1
	userEmail := "test@example.com"

	generatedURL := GenerateSignedURL(userID, userEmail)

	fmt.Println(generatedURL)
	if !strings.Contains(generatedURL, "/verify/email?id=1&user=") {
		t.Errorf("Generated URL does not contain the expected prefix")
	}

	os.Unsetenv("SIGNED_URL_SECRET")
}
