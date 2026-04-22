package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GenerateSessionToken() (plainToken string, tokenHash string, err error) {
	buffer := make([]byte, 32)
	if _, err = rand.Read(buffer); err != nil {
		return "", "", err
	}

	plainToken = base64.RawURLEncoding.EncodeToString(buffer)
	sum := sha256.Sum256([]byte(plainToken))
	tokenHash = hex.EncodeToString(sum[:])

	return plainToken, tokenHash, nil
}
