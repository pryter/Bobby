package pushEvent

import (
	"bobby-worker/internal/cli"
	"path"
	"strconv"
)

// LocalPathVariables Setup path variables

type PathVarSetupOptions struct {
	BuildOutputFolder string
	RuntimeRoot       string
}

func SetupPathVars(
	repositoryID int64,
	options PathVarSetupOptions,
) cli.LocalPathVariables {

	lockerPath := path.Join(
		options.RuntimeRoot, strconv.FormatInt(repositoryID, 10),
	)

	repoPath := path.Join(lockerPath, "repo")
	artifactsOutPath := path.Join(lockerPath, "artifacts")

	artifactSource := path.Join(repoPath, ".next")

	if options.BuildOutputFolder != "" {
		artifactSource = path.Join(repoPath, options.BuildOutputFolder)
	}

	return cli.LocalPathVariables{
		Locker:         lockerPath,
		Repository:     repoPath,
		ArtifactOut:    artifactsOutPath,
		ArtifactSource: artifactSource,
	}
}

// Setup build environments

func SetupBuildEnvironment(
	preset string,
	options cli.BuildEnvironment,
) cli.BuildEnvironment {
	switch preset {
	case "node-default":
		return cli.BuildEnvironment{
			InitCommand:  cli.ExecutableCommand{"yarn", ""},
			BuildCommand: cli.ExecutableCommand{"yarn", "build"},
			EnvType:      cli.EnvNode,
		}
	}

	return cli.BuildEnvironment{}
}
