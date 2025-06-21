package deepl

import (
	"os"

	"github.com/bounoable/deepl"
)

// NewClient creates a new DeepL client.
// It reads the authentication key from the SYSTEM_DEEPL_API_KEY environment variable.
func NewClient() *deepl.Client {
	authKey := os.Getenv("SYSTEM_DEEPL_API_KEY")
	if authKey == "" {
		return nil
	}

	client := deepl.New(
		authKey,
		deepl.BaseURL(os.Getenv("SYSTEM_DEEPL_BASE_URL")),
	)

	return client
}
