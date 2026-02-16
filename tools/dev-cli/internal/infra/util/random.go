package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func RandomSecret(byteLen int) (string, error) {
	buf := make([]byte, byteLen)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate random secret: %w", err)
	}
	return base64.RawStdEncoding.EncodeToString(buf), nil
}
