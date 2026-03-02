package security

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EncryptionKey() []byte {
	// Optional for local development; production should inject env directly.
	_ = godotenv.Load(".env")

	keyBase64 := os.Getenv("ENCRYPTION_KEY")
	if keyBase64 == "" {
		log.Fatal("ENCRYPTION_KEY not set")
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		log.Fatal("invalid ENCRYPTION_KEY format")
	}

	return key
}
