package main

import (
	"bobby-worker/cmd"
	"bobby-worker/internal/utils"
	"bytes"
	_ "embed"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
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
	AppVersion   string `mapstructure:"app_version"`
	HTTPServices struct {
		Worker    cmd.HTTPServiceConfig `mapstructure:"worker"`
		Artifacts cmd.HTTPServiceConfig `mapstructure:"artifacts_server"`
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
		execPath = utils.GetProjectRoot()
	} else {
		//// Production
		//e, err := os.Executable()
		//if err != nil {
		//	log.Panic().Err(err).Msg("Unable to initialise production root path.")
		//}
		//execPath = path.Dir(e)

		// Platform Specific path

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

		_ = os.Mkdir(appPath, 755)

		execPath = appPath
	}

	replaceVariable(&Configs.HTTPServices.Worker.RuntimeBasePath, "$EXEC_PATH", execPath)
	replaceVariable(&Configs.HTTPServices.Artifacts.RuntimeBasePath, "$EXEC_PATH", execPath)
}
