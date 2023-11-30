package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

// GenerateJWT method generates JWT from APP_TOKEN env variable which
// will be used in GitHub APIs requests.
// Note: This must be used in env loaded environment.
func GenerateJWT() (string, error) {

	fixedToken := strings.ReplaceAll(os.Getenv("APP_TOKEN"), "\\n", "\n")
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(fixedToken))

	if err != nil {
		log.Fatal().Msg("Unable to parse RSA key from PEM.")
		return "", err
	}

	t := jwt.NewWithClaims(
		jwt.SigningMethodRS256, jwt.MapClaims{
			"iat": time.Now().Unix() - 60,
			"exp": time.Now().Unix() + (4 * 60), // expires after 4 minutes
			"iss": os.Getenv("APP_ID"),
		},
	)
	jwtToken, err := t.SignedString(key)

	if err != nil {
		log.Fatal().Msg("Unable to generate JWT token for Git APIs.")
	}

	return jwtToken, nil
}
