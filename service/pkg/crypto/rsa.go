package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/rs/zerolog/log"
)

// CreateRSAKeyPair returns privateKey, publickey and error.
func CreateRSAKeyPair() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Error().Err(err).Msg("Error generating RSA private key")
		return nil, nil, err
	}

	// Encode the private key to the PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Extract the public key from the private key
	publicKey := &privateKey.PublicKey

	// Encode the public key to the PEM format
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}

	return pem.EncodeToMemory(privateKeyPEM), pem.EncodeToMemory(publicKeyPEM), nil
}

func RSADecrypt(cipherText string, privKey []byte) (string, error) {
	block, _ := pem.Decode(privKey)
	if block == nil {
		return "", errors.New("invalid key")
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		return "", errors.New("invalid key")
	}

	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return "", err
	}

	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, priv, ct, label)

	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func RSAEncrypt(secretMessage string, pubKey []byte) (string, error) {

	block, _ := pem.Decode(pubKey)
	if block == nil {
		return "", errors.New("invalid key")
	}
	if got, want := block.Type, "RSA PUBLIC KEY"; got != want {
		return "", errors.New("invalid key")
	}

	// Decode the RSA public key
	key, err := x509.ParsePKCS1PublicKey(block.Bytes)

	if err != nil {
		return "", err
	}

	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, key, []byte(secretMessage), label)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
