package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
func GenerateToken(userID string, roles []string, expiry string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
		"exp":     time.Now().Add(parseDuration(expiry)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
	
}

func ParseToken(tokenString string) (*jwt.MapClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*jwt.MapClaims), nil
}