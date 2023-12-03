package pushEvent

import (
	"fmt"
	"path"
	"strconv"
)

// Setup path variables

type LocalPathVariables struct {
	Locker         string
	Repository     string
	ArtifactSource string
	ArtifactOut    string
}

type PathVarSetupOptions struct {
	BuildOutputFolder string
	RuntimeRoot       string
}

func SetupPathVars(
	repositoryID int64,
	options PathVarSetupOptions,
) LocalPathVariables {

	lockerPath := path.Join(
		options.RuntimeRoot, strconv.FormatInt(repositoryID, 10),
	)

	repoPath := path.Join(lockerPath, "repo")
	artifactsOutPath := path.Join(lockerPath, "artifacts")

	artifactSource := path.Join(repoPath, ".next")

	if options.BuildOutputFolder != "" {
		artifactSource = path.Join(repoPath, options.BuildOutputFolder)
	}

	return LocalPathVariables{
		Locker:         lockerPath,
		Repository:     repoPath,
		ArtifactOut:    artifactsOutPath,
		ArtifactSource: artifactSource,
	}
}

// Setup build environments

type EnvironmentType string

const (
	EnvNode EnvironmentType = "node"
)

type ExecutableCommand struct {
	Name string
	Args string
}

func (r ExecutableCommand) flattened() string {
	return fmt.Sprintf("%s %s", r.Name, r.Args)
}

type BuildEnvironment struct {
	InitCommand  ExecutableCommand
	BuildCommand ExecutableCommand
	EnvType      EnvironmentType
}

func SetupBuildEnvironment(
	preset string,
	options BuildEnvironment,
) BuildEnvironment {
	switch preset {
	case "node-default":
		return BuildEnvironment{
			InitCommand:  ExecutableCommand{"yarn", ""},
			BuildCommand: ExecutableCommand{"yarn", "build"},
			EnvType:      EnvNode,
		}
	}

	return BuildEnvironment{}
}
