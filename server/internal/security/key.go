package security

import (
	"encoding/base64"
	"log"
	"os"
	"server/internal/envloader"
	"strings"
)

func EncryptionKey() []byte {
	envloader.LoadDotEnv()

	keyBase64 := strings.TrimSpace(os.Getenv("ENCRYPTION_KEY"))
	if keyBase64 == "" {
		log.Fatalf("[security] ENCRYPTION_KEY is missing or empty. Set ENCRYPTION_KEY in your environment or .env file.")
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		log.Fatalf("[security] invalid ENCRYPTION_KEY format")
	}

	return key
}
