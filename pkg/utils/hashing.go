package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

var secret = []byte("server-secret")

func Hashing(token string) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
