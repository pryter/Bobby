package pushEvent

import (
	"Bobby/internal/utils"
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
}

func SetupPathVars(repositoryID int64, options PathVarSetupOptions) LocalPathVariables {

	prjRoot := utils.GetProjectRoot()
	lockerPath := path.Join(prjRoot, "locker", strconv.FormatInt(repositoryID, 10))

	repoPath := path.Join(lockerPath, "repo")
	artifactsOutPath := path.Join(lockerPath, "artifacts")

	artifactSource := path.Join(repoPath, ".next")

	if options.BuildOutputFolder != "" {
		artifactSource = path.Join(repoPath, options.BuildOutputFolder)
	}

	return LocalPathVariables{Locker: lockerPath, Repository: repoPath, ArtifactOut: artifactsOutPath, ArtifactSource: artifactSource}
}

// Setup build environments

type BuildEnvSetupOptions struct {
	initCMD  string
	buildCMD string
}

func SetupBuildEnvironment(preset string, options BuildEnvironment) BuildEnvironment {
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
