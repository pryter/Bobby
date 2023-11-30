package pushEvent

import (
	"Bobby/internal/gitAPI"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/rs/zerolog/log"
	"strconv"
)

type logFactory struct {
	NoCMD func(err error, msg string)
	CMD   func(err error, msg string, cmd ExecutableCommand)
}

func newLogFactory(repoID int64) logFactory {
	return logFactory{
		func(err error, msg string) {
			log.Fatal().Err(err).Str(
				"repo_id", strconv.Itoa(int(repoID)),
			).Msg(msg)
		}, func(err error, msg string, cmd ExecutableCommand) {
			log.Fatal().Err(err).Str(
				"repo_id", strconv.Itoa(int(repoID)),
			).Str("command", cmd.flattened()).Msg(msg)
		},
	}
}

func WebhookPushEvent(payload github.PushPayload) {

	/*
		TODO
		[X] generate app access token
		[X] clone or pull repository
		[X] build project
		[X] create commit check run
		[ ] provide artifacts url
		[ ] log tunnel
		[ ] error handling
	*/

	// Init required variables
	installID := payload.Installation.ID
	repoID := payload.Repository.ID

	lf := newLogFactory(repoID)

	// Issue github's access token
	accessToken, _ := gitAPI.IssueAccessToken(installID, repoID)
	checkrun := gitAPI.NewCheckRun(
		payload.Repository.HooksURL, payload.Commits[0].ID, accessToken,
	)

	// Initiate CLI tools
	cli := cliFactory{
		pathVars: SetupPathVars(repoID, PathVarSetupOptions{}),
		gitToken: accessToken,
		buildEnv: SetupBuildEnvironment("node-default", BuildEnvironment{}),
	}

	// Clone or pull repository from remote source
	err := cli.CloneRepoWithToken(payload.Repository.CloneURL)

	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		if err := cli.PullChanges(); err != nil {
			lf.NoCMD(err, "Can not pull changes from the remote.")

			checkrun.Update(
				"completed", "failed",
				gitAPI.CheckRunOutput{
					Title:   "Repository Error",
					Summary: "Build server can not pull latest changes from this repo.",
				},
			)
			return
		}
	} else if err != nil {
		return
	}

	// Project workflow
	// init project
	if err := cli.InitProject(); err != nil {
		lf.CMD(
			err, "Unable to initialise the project.", cli.buildEnv.InitCommand,
		)

		checkrun.Update(
			"completed", "failed",
			gitAPI.CheckRunOutput{
				Title:   "Project Error",
				Summary: "Build server can not initialise the project.",
			},
		)
		return
	}

	// build project
	if err := cli.Build(); err != nil {
		lf.CMD(
			err, "Unable to build the project.", cli.buildEnv.BuildCommand,
		)

		checkrun.Update(
			"completed", "failed",
			gitAPI.CheckRunOutput{
				Title:   "Build Error",
				Summary: "Build server can not build the project.",
			},
		)
		return
	}

	// Export and compress artifacts to zip file
	cli.ExportArtifact()

	checkrun.Update(
		"completed", "success", gitAPI.CheckRunOutput{
			Title: "Build Success",
			Summary: "Successfully built the changes. " +
				"The artifact file can be access at https://bobby.pryter.me/some_artifact_uri",
		},
	)

}
