package pushEvent

import (
	"Bobby/internal/utils"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type cliFactory struct {
	pathVars LocalPathVariables
	buildEnv BuildEnvironment
	gitToken string
}

func (r *cliFactory) CloneRepoWithToken(cloneURL string) error {
	authCloneURl := strings.ReplaceAll(cloneURL, "github.com", fmt.Sprintf("x-access-token:%s@github.com", r.gitToken))

	_, err := git.PlainClone(r.pathVars.Repository, false, &git.CloneOptions{
		URL:      authCloneURl,
		Progress: os.Stdout,
	})

	return err
}

func (r *cliFactory) PullChanges() error {
	activeRepo, err := git.PlainOpen(r.pathVars.Repository)
	w, err := activeRepo.Worktree()
	if err != nil {
		panic(err)
	}

	err = w.Pull(&git.PullOptions{Auth: &http.BasicAuth{Username: "x-access-token", Password: r.gitToken}})

	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}
	return err
}

func (r *cliFactory) runCommand(command ExecutableCommand) error {
	idepsCMD := exec.Command(command.Name, command.Args)
	idepsCMD.Dir = r.pathVars.Repository
	idepsCMD.Stdout = os.Stdout
	err := idepsCMD.Run()
	return err
}

func (r *cliFactory) InitProject() error {
	return r.runCommand(r.buildEnv.InitCommand)
}

func (r *cliFactory) Build() error {
	return r.runCommand(r.buildEnv.BuildCommand)
}

func (r *cliFactory) ExportArtifact() error {
	artifactFile := fmt.Sprintf("%s/artifact-%d.zip", r.pathVars.ArtifactOut, time.Now().Unix())
	err := os.MkdirAll(r.pathVars.ArtifactOut, os.ModePerm)
	file, _ := os.Create(artifactFile)

	zipIO, _ := utils.ZipDirectory(r.pathVars.ArtifactSource)
	_, err = io.Copy(file, zipIO)
	return err
}
