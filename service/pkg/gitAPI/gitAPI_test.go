package gitAPI

import (
	"Bobby/pkg/utils"
	"github.com/joho/godotenv"
	"path/filepath"
	"testing"
)

func TestAPIS(t *testing.T) {
	godotenv.Load(filepath.Join(utils.GetProjectRoot(), ".env"))
}
