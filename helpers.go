package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func GenerateSessionID(length int) (string, error) {
	// Create a byte slice with the given length
	bytes := make([]byte, length)

	// Fill the slice with random bytes
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	// Convert the byte slice to a hexadecimal string
	return hex.EncodeToString(bytes), nil
}
