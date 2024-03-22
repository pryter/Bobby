package challenge

import "errors"

var (
	ErrPublicKeyMissing      = errors.New("unable to open public key file")
	ErrUnableToGetMacAddress = errors.New("unable to get mac address")
	ErrRSAEncryption         = errors.New("unable to perform RSA encryption to the generated challenge")

	ErrPrivateKeyMissing        = errors.New("unable to open private key file")
	ErrSecretFileMissing        = errors.New("unable to open secret file")
	ErrUnsolvableRSAError       = errors.New("unable to solve the given challenge. (encryption error)")
	ErrUnsolvableIDMismatch     = errors.New("unable to solve the given challenge. (ID mismatch)")
	ErrUnsolvableSecretMismatch = errors.New("unable to solve the given challenge. (secret mismatch)")
)
