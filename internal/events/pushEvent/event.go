package pushEvent

import (
	"Bobby/internal/token"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-playground/webhooks/v6/github"
)

func WebhookPushEvent(payload github.PushPayload) {

	/*
		TODO
		[X] generate app access token
		[X] clone or pull repository
		[X] build project
		[ ] create commit check run
		[ ] provide artifacts url
	*/

	// Init required variables
	installID := payload.Installation.ID
	repoID := payload.Repository.ID

	// Issue github's access token
	accessToken, _ := token.IssueToken(installID, repoID)

	// Initiate CLI tools
	cli := cliFactory{
		pathVars: SetupPathVars(repoID, PathVarSetupOptions{}),
		gitToken: accessToken,
		buildEnv: SetupBuildEnvironment("node-default", BuildEnvironment{}),
	}

	// Clone or pull repository from remote source
	err := cli.CloneRepoWithToken(payload.Repository.CloneURL)

	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		cli.PullChanges()
	}

	// Project workflow
	cli.InitProject()
	cli.Build()

	// Export and compress artifacts to zip file
	cli.ExportArtifact()

}
