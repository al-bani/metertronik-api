package token

import "github.com/google/uuid"

func GenerateRefreshToken() string {
	return uuid.New().String()
}