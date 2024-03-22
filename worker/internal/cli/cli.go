package cli

import (
	"Bobby/pkg/utils"
	"bobby-worker/internal/events/pushEvent"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"
	"time"
)

// CliFactory can be created and fill the field.
// This will provide required cli methods exclusively for pushEvent package.
type CliFactory struct {
	PathVars pushEvent.LocalPathVariables
	BuildEnv pushEvent.BuildEnvironment
}

// CloneRepoWithToken clones a given repo to cli's local machine.
func (r *CliFactory) CloneRepoWithToken(cloneURL string) error {

	auth, err := ssh.DefaultAuthBuilder("git")

	_, err = git.PlainClone(
		r.PathVars.Repository, false, &git.CloneOptions{
			Auth:     auth,
			URL:      cloneURL,
			Progress: os.Stdout,
		},
	)

	log.Print(err)
	return err
}

// PullChanges pulls changes to existed local repository.
func (r *CliFactory) PullChanges() error {
	activeRepo, err := git.PlainOpen(r.PathVars.Repository)
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
		&git.PullOptions{},
	)

	// skip if local repository is already up-to-date.
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}

	return err
}

// runCommand runs ExecutableCommand and controls its output.
func (r *CliFactory) runCommand(command pushEvent.ExecutableCommand) error {
	idepsCMD := exec.Command(command.Name, command.Args)
	idepsCMD.Dir = r.PathVars.Repository
	idepsCMD.Stdout = os.Stdout
	err := idepsCMD.Run()
	return err
}

// InitProject inits the project.
func (r *CliFactory) InitProject() error {
	return r.runCommand(r.BuildEnv.InitCommand)
}

// Build builds the project.
func (r *CliFactory) Build() error {
	return r.runCommand(r.BuildEnv.BuildCommand)
}

// ExportArtifact exports build folder as an artifact zip file.
func (r *CliFactory) ExportArtifact() error {
	artifactFile := fmt.Sprintf(
		"%s/artifact-%d.zip", r.PathVars.ArtifactOut, time.Now().Unix(),
	)

	if err := os.MkdirAll(r.PathVars.ArtifactOut, os.ModePerm); err != nil {
		log.Error().Err(err).Str(
			"path", r.PathVars.ArtifactOut,
		).Msg("Unable to create artifact directory.")
	}

	file, err := os.Create(artifactFile)

	if err != nil {
		log.Error().Err(err).Msg("Unable to create artifact file.")
		return err
	}

	zipIO, err := utils.ZipDirectory(r.PathVars.ArtifactSource)

	if err != nil {
		log.Error().Err(err).Str(
			"path", r.PathVars.ArtifactSource,
		).Msg("Unable zip source directory")
		return err
	}

	_, err = io.Copy(file, zipIO)

	return err
}
