package pushEvent

import (
	"Bobby/internal/token"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"testing"
)

func testEachFactory(cli cliFactory, t *testing.T) {
	t.Run("clone-repository", func(t *testing.T) {
		err := cli.CloneRepoWithToken("https://github.com/pryter/anri.git")
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) && err != nil {
			t.Error(err)
		}
	})

	t.Run("pull-repository", func(t *testing.T) {
		err := cli.PullChanges()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("init-project", func(t *testing.T) {
		err := cli.InitProject()

		if err != nil {
			t.Error(err)
		}
	})

	t.Run("build-project", func(t *testing.T) {
		err := cli.Build()

		if err != nil {
			t.Error(err)
		}
	})

	t.Run("export-artifact", func(t *testing.T) {
		err := cli.ExportArtifact()

		if err != nil {
			t.Error(err)
		}
	})
}

func TestCliFactory(t *testing.T) {

	accessToken, _ := token.IssueToken(44151598, 571145096)

	var clis []cliFactory

	clis = append(clis, cliFactory{
		pathVars: SetupPathVars(571145096, PathVarSetupOptions{}),
		gitToken: accessToken,
		buildEnv: SetupBuildEnvironment("node-default", BuildEnvironment{}),
	})

	for _, s := range clis {
		t.Run(fmt.Sprintf("env-type-%s", s.buildEnv.EnvType), func(st *testing.T) {
			testEachFactory(s, st)
		})
	}

}
