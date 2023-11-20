package token

import (
	"Bobby/internal/utils"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"testing"
)

func TestIssueToken(t *testing.T) {
	t.Run("access-environment-secrets", func(t *testing.T) {

		err := godotenv.Load(filepath.Join(utils.GetProjectRoot(), ".env"))
		if err != nil {
			t.Error(err)
			return
		}

		_, exist := os.LookupEnv("APP_ID")
		if !exist {
			t.Fail()
		}
	})

	t.Run("generate-jwt-token", func(t *testing.T) {
		err := godotenv.Load(filepath.Join(utils.GetProjectRoot(), ".env"))
		if err != nil {
			t.Error(err)
			return
		}

		_, err = generateJWT()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("issue-access-token", func(t *testing.T) {
		err := godotenv.Load(filepath.Join(utils.GetProjectRoot(), ".env"))
		if err != nil {
			t.Error(err)
			return
		}

		token, err := IssueToken(44151598, 571145096)

		if err != nil {
			t.Error(err)
		}

		t.Logf("TOKEN: %s", token)
	})
}
