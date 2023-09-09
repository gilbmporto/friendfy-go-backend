package auth

import (
	"friendfy-api/src/config"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CreateToken(userID uint64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["userID"] = userID

	// secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	strToken, err := token.SignedString([]byte(config.SecretKey))
	return strToken, err
}
