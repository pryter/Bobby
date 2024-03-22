package challenge

import (
	"Bobby/pkg/crypto"
	"bobby-worker/internal/utils"
	"os"
	"time"
)

type Builder struct {
	PublicKeyPath string
	SetupID       string
}

// Generate challenge
func (b Builder) Generate() (string, error) {

	pubKey, err := os.ReadFile(b.PublicKeyPath)

	if err != nil {
		return "", ErrPublicKeyMissing
	}

	macAddr, err := utils.GetMacAddr()
	if err != nil || macAddr[0] == "" {
		return "", ErrUnableToGetMacAddress
	}

	challengeRawText := macAddr[0] + "|" + b.SetupID + "|" + time.Now().String()
	challenge, err := crypto.RSAEncrypt(challengeRawText, pubKey)

	if err != nil {
		return "", ErrRSAEncryption
	}

	return challenge, nil
}
