package challenge

import (
	"Bobby/pkg/crypto"
	"net"
	"os"
	"time"
)

type Builder struct {
	PublicKeyPath string
	SetupID       string
}

func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}

	return as, nil
}

// Generate challenge
func (b Builder) Generate() (string, error) {

	pubKey, err := os.ReadFile(b.PublicKeyPath)

	if err != nil {
		return "", ErrPublicKeyMissing
	}

	macAddr, err := getMacAddr()
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
