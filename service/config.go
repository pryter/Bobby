package main

import (
	"Bobby/cmd"
	"Bobby/pkg/utils"
	"bytes"
	_ "embed"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path"
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
		Webhook   cmd.HTTPServiceConfig `mapstructure:"webhook"`
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

	var basePath string
	if strings.Contains(os.Args[0], "/tmp/") {
		basePath = utils.GetProjectRoot()
	} else {
		e, err := os.Executable()
		if err != nil {
			log.Panic().Err(err).Msg("Unable to initialise production root path.")
		}
		basePath = path.Dir(e)
	}

	replaceVariable(&Configs.HTTPServices.Webhook.RuntimeBasePath, "$EXEC_PATH", basePath)
	replaceVariable(&Configs.HTTPServices.Artifacts.RuntimeBasePath, "$EXEC_PATH", basePath)
}
