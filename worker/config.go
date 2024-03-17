package main

import (
	"Bobby/pkg/utils"
	"bobby-worker/cmd"
	"bytes"
	_ "embed"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path"
	"runtime"
	"strings"
)

//go:embed build-config.yaml
var configFile []byte

var Configs Config

type ConcurrentPoolConfig struct {
	MaxConcurrentTasks int `mapstructure:"max_concurrent_tasks"`
}

type Config struct {
	AppVersion      string `mapstructure:"app_version"`
	AppResourcePath string `mapstructure:"app_resource_path"`
	HTTPServices    struct {
		Worker    cmd.WorkerServiceOptions   `mapstructure:"worker"`
		Artifacts cmd.ArtifactServiceOptions `mapstructure:"artifacts_server"`
	} `mapstructure:"http_services"`
	ConcurrentPool ConcurrentPoolConfig `mapstructure:"concurrent_pool"`
}

func replaceVariable(bp *string, varName string, varVal string) {
	*bp = strings.ReplaceAll(
		*bp, varName, varVal,
	)
}

func init() {

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(configFile))

	if err != nil {
		log.Panic().Err(err).Msg("Unable to find bobby-notes")
	}

	err = viper.Unmarshal(&Configs)

	if err != nil {
		log.Panic().Err(err).Msg("Unable to unpack config file.")
	}

	var execPath string
	if strings.Contains(os.Args[0], "/tmp/") {
		// Development
		execPath = path.Join(utils.GetProjectRoot(), "resources/Bobby-worker")
		_ = os.Mkdir(execPath, 0777)
	} else {
		//// Production
		var appPath string
		switch runtime.GOOS {
		case "windows":
			appPath = "C:\\ProgramData\\Bobby-worker"
			break
		case "darwin":
			appPath = "/Library/Application Support/Bobby-worker"
			break
		case "linux":
			appPath = "/var/lib/Bobby-worker"
			break
		default:
			log.Panic().Err(err).Msg("Unable to identify OS.")
		}

		_ = os.Mkdir(appPath, 0777)

		execPath = appPath
	}

	replaceVariable(&Configs.AppResourcePath, "$RESOURCE_PATH", execPath)
	replaceVariable(&Configs.HTTPServices.Worker.ServiceBasePath, "$RESOURCE_PATH", execPath)
	replaceVariable(&Configs.HTTPServices.Artifacts.RuntimeBasePath, "$RESOURCE_PATH", execPath)
}
