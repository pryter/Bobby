package token

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
	"time"
)

func GenerateJWT() (string, error) {

	fixedToken := strings.ReplaceAll(os.Getenv("APP_TOKEN"), "\\n", "\n")
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(fixedToken))

	if err != nil {
		panic(err)
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
