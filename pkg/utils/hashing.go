package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"metertronik/pkg/config"
)


func Hashing(token string) string {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	secret := []byte(cfg.SECRETKEY)

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
