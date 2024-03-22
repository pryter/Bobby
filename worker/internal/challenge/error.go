package challenge

import "errors"

var (
	ErrPublicKeyMissing      = errors.New("unable to open public key file")
	ErrUnableToGetMacAddress = errors.New("unable to get mac address")
	ErrRSAEncryption         = errors.New("unable to perform RSA encryption to the generated challenge")
)
