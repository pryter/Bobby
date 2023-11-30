package pushEvent

import (
	"Bobby/internal/utils"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// cliFactory can be created and fill the field.
// This will provide required cli methods exclusively for pushEvent package.
type cliFactory struct {
	pathVars LocalPathVariables
	buildEnv BuildEnvironment
	gitToken string
}

// CloneRepoWithToken clones a given repo to cli's local machine.
func (r *cliFactory) CloneRepoWithToken(cloneURL string) error {
	authCloneURl := strings.ReplaceAll(
		cloneURL, "github.com",
		fmt.Sprintf("x-access-token:%s@github.com", r.gitToken),
	)

	_, err := git.PlainClone(
		r.pathVars.Repository, false, &git.CloneOptions{
			URL:      authCloneURl,
			Progress: os.Stdout,
		},
	)

	return err
}

// PullChanges pulls changes to existed local repository.
func (r *cliFactory) PullChanges() error {
	activeRepo, err := git.PlainOpen(r.pathVars.Repository)
	if err != nil {
		log.Error().Err(err).Msg("Unable to open local repository.")
		return err
	}

	w, err := activeRepo.Worktree()
	if err != nil {
		log.Error().Err(err).Msg("Unable get repository's worktree.")
		return err
	}

	err = w.Pull(
		&git.PullOptions{
			Auth: &http.BasicAuth{
				Username: "x-access-token",
				Password: r.gitToken,
			},
		},
	)

	// skip if local repository is already up-to-date.
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}

	return err
}

// runCommand runs ExecutableCommand and controls its output.
func (r *cliFactory) runCommand(command ExecutableCommand) error {
	idepsCMD := exec.Command(command.Name, command.Args)
	idepsCMD.Dir = r.pathVars.Repository
	idepsCMD.Stdout = os.Stdout
	err := idepsCMD.Run()
	return err
}

// InitProject inits the project.
func (r *cliFactory) InitProject() error {
	return r.runCommand(r.buildEnv.InitCommand)
}

// Build builds the project.
func (r *cliFactory) Build() error {
	return r.runCommand(r.buildEnv.BuildCommand)
}

// ExportArtifact exports build folder as an artifact zip file.
func (r *cliFactory) ExportArtifact() error {
	artifactFile := fmt.Sprintf(
		"%s/artifact-%d.zip", r.pathVars.ArtifactOut, time.Now().Unix(),
	)

	if err := os.MkdirAll(r.pathVars.ArtifactOut, os.ModePerm); err != nil {
		log.Error().Err(err).Str(
			"path", r.pathVars.ArtifactOut,
		).Msg("Unable to create artifact directory.")
	}

	file, err := os.Create(artifactFile)

	if err != nil {
		log.Error().Err(err).Msg("Unable to create artifact file.")
		return err
	}

	zipIO, err := utils.ZipDirectory(r.pathVars.ArtifactSource)

	if err != nil {
		log.Error().Err(err).Str(
			"path", r.pathVars.ArtifactSource,
		).Msg("Unable zip source directory")
		return err
	}

	_, err = io.Copy(file, zipIO)

	return err
}
