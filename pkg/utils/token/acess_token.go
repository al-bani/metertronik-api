package token

import (
	"metertronik/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("metertronik")

func GenerateAccessToken(userID int64) string {
	claims := jwt.MapClaims{
		"uid": userID,
		"exp": utils.TimeNow().Add(utils.Minutes(15)).Time.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(jwtSecret)
	return signed
}
