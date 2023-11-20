package pushEvent

import (
	"Bobby/internal/utils"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func testSetupPathVarsWithOpts(ROOT string, option PathVarSetupOptions, t *testing.T) {
	pathVars := SetupPathVars(1234567, option)

	a := assert.New(t)

	outputFolder := ".next"
	if option.BuildOutputFolder != "" {
		outputFolder = option.BuildOutputFolder
	}

	a.Equal(pathVars.ArtifactOut, path.Join(ROOT, "locker/1234567/artifacts"), "Path ArtifactOut Mismatch")
	a.Equal(pathVars.ArtifactSource, path.Join(ROOT, "locker/1234567/repo/", outputFolder), "Path ArtifactSource Mismatch")
	a.Equal(pathVars.Locker, path.Join(ROOT, "locker/1234567"), "Path Locker Mismatch")
	a.Equal(pathVars.Repository, path.Join(ROOT, "locker/1234567/repo"), "Path Repo Mismatch")
}

func TestEnvironmentSetup(t *testing.T) {
	t.Run("setup-path-variables", func(t *testing.T) {
		ROOT := utils.GetProjectRoot()
		testSetupPathVarsWithOpts(ROOT, PathVarSetupOptions{}, t)
	})

	t.Run("setup-path-variables-with-options", func(t *testing.T) {
		ROOT := utils.GetProjectRoot()
		testSetupPathVarsWithOpts(ROOT, PathVarSetupOptions{BuildOutputFolder: "somethingelse"}, t)
	})
}
