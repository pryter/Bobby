package pushEvent

type EnvironmentType string

const (
	EnvNode EnvironmentType = "node"
)

type ExecutableCommand struct {
	Name string
	Args string
}

type BuildEnvironment struct {
	InitCommand  ExecutableCommand
	BuildCommand ExecutableCommand
	EnvType      EnvironmentType
}
