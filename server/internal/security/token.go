package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func HashNormalizedEmail(email string) string {
	normalized := strings.ToLower(strings.TrimSpace(email))
	mac := hmac.New(sha256.New, EncryptionKey())
	mac.Write([]byte(normalized))
	return hex.EncodeToString(mac.Sum(nil))
}
