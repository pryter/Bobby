package challenge

import (
	"Bobby/pkg/crypto"
	"os"
	"strings"
)

type Solver struct {
	PrivateKeyPath string
	SecretPath     string
	SetupID        string
}

// Solve method solves given challenge and give the result whether the challenge is properly solvable or not
func (s Solver) Solve(challenge string) (bool, error) {
	key, err := os.ReadFile(s.PrivateKeyPath)
	if err != nil {
		return false, ErrPrivateKeyMissing
	}
	secret, err := os.ReadFile(s.SecretPath)
	if err != nil {
		return false, ErrSecretFileMissing
	}

	decrypted, err := crypto.RSADecrypt(challenge, key)
	if err != nil {
		return false, ErrUnsolvableRSAError
	}

	slice := strings.Split(decrypted, "|")

	if slice[1] != s.SetupID {
		return false, ErrUnsolvableIDMismatch
	}

	if slice[0] != string(secret) {
		return false, ErrUnsolvableSecretMismatch
	}

	return true, nil
}
