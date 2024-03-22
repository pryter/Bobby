package pushEvent

import (
	"Bobby/pkg/utils"
	"bobby-worker/internal/cli"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"path"
	"path/filepath"
	"testing"
)

func testEachFactory(cli cli.CliFactory, t *testing.T) {
	t.Run(
		"clone-repository", func(t *testing.T) {
			err := cli.CloneRepoWithToken("git@github.com:pryter/anri.git")
			if !errors.Is(err, git.ErrRepositoryAlreadyExists) && err != nil {
				t.Error(err)
			}
		},
	)

	t.Run(
		"pull-repository", func(t *testing.T) {
			err := cli.PullChanges()
			if err != nil {
				t.Error(err)
			}
		},
	)

	t.Run(
		"init-project", func(t *testing.T) {
			err := cli.InitProject()

			if err != nil {
				t.Error(err)
			}
		},
	)

	t.Run(
		"build-project", func(t *testing.T) {
			err := cli.Build()

			if err != nil {
				t.Error(err)
			}
		},
	)

	t.Run(
		"export-artifact", func(t *testing.T) {
			err := cli.ExportArtifact()

			if err != nil {
				t.Error(err)
			}
		},
	)
}

func TestEnvironmentSetup(t *testing.T) {
	t.Run(
		"network-path-variables", func(t *testing.T) {
			ROOT := utils.GetProjectRoot()
			testSetupPathVarsWithOpts(ROOT, PathVarSetupOptions{}, t)
		},
	)

	t.Run(
		"network-path-variables-with-options", func(t *testing.T) {
			ROOT := utils.GetProjectRoot()
			testSetupPathVarsWithOpts(
				ROOT, PathVarSetupOptions{BuildOutputFolder: "somethingelse"}, t,
			)
		},
	)
}

func TestCliFactory(t *testing.T) {

	err := godotenv.Load(filepath.Join(utils.GetProjectRoot(), ".env"))
	if err != nil {
		t.Error(err)
		return
	}

	var clis []cli.CliFactory

	clis = append(
		clis, cli.CliFactory{
			PathVars: SetupPathVars(571145096, PathVarSetupOptions{}),
			BuildEnv: SetupBuildEnvironment("node-default", BuildEnvironment{}),
		},
	)

	for _, s := range clis {
		t.Run(
			fmt.Sprintf("env-type-%s", s.BuildEnv.EnvType), func(st *testing.T) {
				testEachFactory(s, st)
			},
		)
	}

}

func testSetupPathVarsWithOpts(ROOT string, option PathVarSetupOptions, t *testing.T) {
	pathVars := SetupPathVars(1234567, option)

	a := assert.New(t)

	outputFolder := ".next"
	if option.BuildOutputFolder != "" {
		outputFolder = option.BuildOutputFolder
	}

	a.Equal(
		pathVars.ArtifactOut, path.Join(ROOT, "locker/1234567/artifacts"),
		"Path ArtifactOut Mismatch",
	)
	a.Equal(
		pathVars.ArtifactSource, path.Join(ROOT, "locker/1234567/repo/", outputFolder),
		"Path ArtifactSource Mismatch",
	)
	a.Equal(pathVars.Locker, path.Join(ROOT, "locker/1234567"), "Path Locker Mismatch")
	a.Equal(pathVars.Repository, path.Join(ROOT, "locker/1234567/repo"), "Path Repo Mismatch")
}

type Commits []struct {
	Sha       string `json:"sha"`
	ID        string `json:"id"`
	NodeID    string `json:"node_id"`
	TreeID    string `json:"tree_id"`
	Distinct  bool   `json:"distinct"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Author    struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"author"`
	Committer struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"committer"`
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}
type Commit struct {
	Sha       string `json:"sha"`
	ID        string `json:"id"`
	NodeID    string `json:"node_id"`
	TreeID    string `json:"tree_id"`
	Distinct  bool   `json:"distinct"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Author    struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"author"`
	Committer struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"committer"`
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}

func TestWebhookPushEvent(t *testing.T) {

	godotenv.Load(filepath.Join(utils.GetProjectRoot(), ".env"))

	repository := github.PushPayload{}.Repository
	repository.ID = 571145096
	repository.CloneURL = "https://github.com/pryter/anri.git"
	repository.HooksURL = "https://api.github.com/repos/pryter/forms/hooks"

	installation := github.PushPayload{}.Installation
	installation.ID = 44151598

	var commits Commits
	commit := Commit{ID: "72a2d245818eec4471816aee9ad4e01dca5d3aa2"}
	commits = append(commits, commit)

	WebhookPushEvent(json.RawMessage{}, WebhookPushEventOptions{})

}
