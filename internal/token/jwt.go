package token

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func generateJWT() (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("APP_KEY")))

	if err != nil {
		return "", err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix() - 60,
		"exp": time.Now().Unix() + (10 * 60),
		"iss": os.Getenv("APP_ID"),
	})
	jwtToken, _ := t.SignedString(key)

	return jwtToken, nil
}
