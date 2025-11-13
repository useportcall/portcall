package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
)

type ICrypto interface {
	Encrypt(data string) (string, error)
	Decrypt(data string) (string, error)
	CompareHash(hashed string, plain string) (bool, error)
}

type crypto struct {
	key []byte
}

func New() ICrypto {
	b64key := os.Getenv("AES_ENCRYPTION_KEY")
	if b64key == "" {
		log.Fatal("AES_ENCRYPTION_KEY environment variable not set")
	}

	key, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		log.Fatalf("failed to encrypt restricted key: %v", err)
	}

	return &crypto{key: key}
}

// Encrypts the plaintext API key using AES-GCM
func (c *crypto) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize()) // Nonce must be unique for each encryption
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *crypto) Decrypt(encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)

	return string(plaintext), err
}

func (c *crypto) CompareHash(hashed string, plain string) (bool, error) {
	decrypted, err := c.Decrypt(hashed)
	if err != nil {
		return false, err
	}

	return decrypted == plain, nil
}
