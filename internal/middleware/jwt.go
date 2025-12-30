package middleware

import (
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
)

var jwtSecret = []byte("metertronik")

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("JWTMiddleware")
		auth := c.GetHeader("Authorization")
		
		if !strings.HasPrefix(auth, "Bearer ") {
			log.Println("Token Not Valid: missing token")
			c.JSON(401, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("Token Not Valid: ")
			c.JSON(401, gin.H{"error": "token expired"})
			c.Abort()
			return
		}

		log.Println("Token Valid? : ", token.Valid)

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", int(claims["uid"].(float64)))

		log.Println("passing middleware")

		c.Next()
	}
}
