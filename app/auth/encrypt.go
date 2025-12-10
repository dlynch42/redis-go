package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func EncryptPassword(password string) string {
	// Strip ">" prefix if present
	if len(password) > 0 && password[0] == '>' {
		password = password[1:]
	}

	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
