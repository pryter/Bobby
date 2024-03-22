package pushEvent

import (
	"Bobby/pkg/comm"
	"bobby-worker/internal/cli"
	"bobby-worker/internal/hostAPI"
	"encoding/json"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"strconv"
)

type logFactory struct {
	NoCMD func(err error, msg string)
	CMD   func(err error, msg string, cmd cli.ExecutableCommand)
}

func newLogFactory(repoID int64) logFactory {
	return logFactory{
		func(err error, msg string) {
			log.Error().Err(err).Str(
				"repo_id", strconv.Itoa(int(repoID)),
			).Msg(msg)
		}, func(err error, msg string, cmd cli.ExecutableCommand) {
			log.Warn().Err(err).Str(
				"repo_id", strconv.Itoa(int(repoID)),
			).Str("command", cmd.Flattened()).Msg(msg)
		},
	}
}

type WebhookPushEventOptions struct {
	RuntimeBasePath string
}

func WebhookPushEvent(
	rawPayload json.RawMessage,
	conn *websocket.Conn,
	options WebhookPushEventOptions,
) {

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

	var payload github.PushPayload
	err := json.Unmarshal(rawPayload, &payload)
	if err != nil {

	}

	// Init required variables
	//installID := payload.Installation.ID
	repoID := payload.Repository.ID

	lf := newLogFactory(repoID)

	checkRunBody := comm.CheckRunBody{
		Name:    "Bobby",
		HeadSHA: payload.Commits[0].ID,
		Status:  "in_progress",
		Output: comm.CheckRunOutput{
			Title: "Building in progress",
			Summary: "Build server is building your project" +
				"\nFor more information visit https://bobby.pryter.me/task_id/log",
		},
	}

	checkrun, err := hostAPI.NewCheckRunAPI(
		conn, checkRunBody, payload.Repository.HooksURL, payload.Repository.ID,
		payload.Installation.ID,
	)

	// Initiate CLI tools
	cliTool := cli.CliFactory{
		PathVars: SetupPathVars(repoID, PathVarSetupOptions{RuntimeRoot: options.RuntimeBasePath}),
		BuildEnv: SetupBuildEnvironment("node-default", cli.BuildEnvironment{}),
	}

	// Clone or pull repository from remote source
	err = cliTool.CloneRepoWithToken(payload.Repository.SSHURL)

	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		if err := cliTool.PullChanges(); err != nil {
			lf.NoCMD(err, "Can not pull changes from the remote.")

			checkrun.Update(
				"completed", comm.ConclusionFailure,
				comm.CheckRunOutput{
					Title:   "Repository Error",
					Summary: "Build server can not pull latest changes from this repo.",
				},
			)
			return
		}
	} else if err != nil {
		checkrun.Update(
			"completed", comm.ConclusionFailure,
			comm.CheckRunOutput{
				Title:   "Repository Error",
				Summary: "An unexpected error occurs on the worker unit",
			},
		)
		lf.NoCMD(err, "unable to clone")
		return
	}

	// Project workflow
	// init project
	if err := cliTool.InitProject(); err != nil {
		lf.CMD(
			err, "Unable to initialise the project.", cliTool.BuildEnv.InitCommand,
		)

		checkrun.Update(
			"completed", comm.ConclusionFailure,
			comm.CheckRunOutput{
				Title:   "Project Error",
				Summary: "Build server can not initialise the project.",
			},
		)
		return
	}

	// build project
	if err := cliTool.Build(); err != nil {
		lf.CMD(
			err, "Unable to build the project.", cliTool.BuildEnv.BuildCommand,
		)

		checkrun.Update(
			"completed", comm.ConclusionFailure,
			comm.CheckRunOutput{
				Title:   "Build Error",
				Summary: "Build server can not build the project.",
			},
		)
		return
	}

	// Export and compress artifacts to zip file
	cliTool.ExportArtifact()

	checkrun.Update(
		"completed", comm.ConclusionSuccess, comm.CheckRunOutput{
			Title: "Build Success",
			Summary: "Successfully built the changes. " +
				"The artifact file can be access at https://bobby.pryter.me/some_artifact_uri",
		},
	)

}
