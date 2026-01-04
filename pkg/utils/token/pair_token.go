package token

import (
	"crypto/rand"
	"encoding/base64"
)

func GeneratePairToken() string {
b := make([]byte, 24)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}